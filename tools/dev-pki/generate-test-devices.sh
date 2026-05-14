#!/usr/bin/env bash
set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd -P)"
out_dir="${OUT_DIR:-$repo_root/keys/test_device}"
device_count="${DEVICE_COUNT:-100}"
device_valid_days="${DEVICE_VALID_DAYS:-180}"
ca_valid_days="${CA_VALID_DAYS:-365}"

if [[ -e "$out_dir" && "${FORCE:-false}" != "true" ]]; then
  echo "$out_dir already exists; set FORCE=true to regenerate" >&2
  exit 2
fi

rm -rf "$out_dir"
mkdir -p "$out_dir/ca" "$out_dir/devices" "$out_dir/bundles"

openssl ecparam -name prime256v1 -genkey -noout -out "$out_dir/ca/dev-device-ca.key"
chmod 0600 "$out_dir/ca/dev-device-ca.key"
cat > "$out_dir/ca/ca.conf" <<EOF
[ req ]
prompt = no
distinguished_name = dn
x509_extensions = v3_ca

[ dn ]
C = TW
O = Realtek Connect Plus Development
OU = Test Device PKI
CN = Realtek Connect Plus Dev Device CA

[ v3_ca ]
basicConstraints = critical, CA:TRUE, pathlen:0
keyUsage = critical, keyCertSign, cRLSign
subjectKeyIdentifier = hash
authorityKeyIdentifier = keyid:always,issuer
EOF

openssl req -x509 -new -sha256 \
  -key "$out_dir/ca/dev-device-ca.key" \
  -days "$ca_valid_days" \
  -out "$out_dir/ca/dev-device-ca.crt" \
  -config "$out_dir/ca/ca.conf"

serial_file="$out_dir/ca/dev-device-ca.srl"
manifest="$out_dir/manifest.csv"
printf 'device_id,product_type,model,display_name,certificate_path,key_path,bundle_path\n' > "$manifest"

json_array_started=false
json_manifest="$out_dir/manifest.json"
printf '[\n' > "$json_manifest"

device_type_for_index() {
  local n="$1"
  if (( n <= 15 )); then echo "light_bulb"
  elif (( n <= 30 )); then echo "air_conditioner"
  elif (( n <= 42 )); then echo "smart_plug"
  elif (( n <= 62 )); then echo "camera"
  elif (( n <= 72 )); then echo "gateway"
  elif (( n <= 84 )); then echo "sensor"
  elif (( n <= 90 )); then echo "switch_panel"
  elif (( n <= 96 )); then echo "curtain"
  else echo "linux_simulator"
  fi
}

model_for_type() {
  case "$1" in
    light_bulb) echo "RTC-LB-RGBTW-01" ;;
    air_conditioner) echo "RTC-AC-IR-01" ;;
    smart_plug) echo "RTC-SP-PM-01" ;;
    camera) echo "RTC-CAM-PRO2-01" ;;
    gateway) echo "RTC-GW-BRIDGE-01" ;;
    sensor) echo "RTC-SEN-THM-01" ;;
    switch_panel) echo "RTC-SW-4GANG-01" ;;
    curtain) echo "RTC-CURTAIN-01" ;;
    linux_simulator) echo "RTC-SIM-LINUX-01" ;;
  esac
}

display_name_for_type() {
  local type="$1"
  local ordinal="$2"
  case "$type" in
    light_bulb) echo "Tunable Light Bulb $ordinal" ;;
    air_conditioner) echo "Air Conditioner $ordinal" ;;
    smart_plug) echo "Metering Smart Plug $ordinal" ;;
    camera) echo "PRO2 Camera Demo $ordinal" ;;
    gateway) echo "Bridge Gateway $ordinal" ;;
    sensor) echo "Environmental Sensor $ordinal" ;;
    switch_panel) echo "Scene Switch Panel $ordinal" ;;
    curtain) echo "Motorized Curtain $ordinal" ;;
    linux_simulator) echo "Linux Simulator $ordinal" ;;
  esac
}

capabilities_for_type() {
  case "$1" in
    light_bulb) echo '["power", "brightness", "color_temperature", "rgb_color", "firmware_ota"]' ;;
    air_conditioner) echo '["power", "mode", "target_temperature", "fan_speed", "swing", "firmware_ota"]' ;;
    smart_plug) echo '["power", "power_meter", "energy_total", "overload_state", "firmware_ota"]' ;;
    camera) echo '["power", "snapshot", "stream_webrtc", "motion_detection", "night_vision", "two_way_audio", "firmware_ota"]' ;;
    gateway) echo '["child_device_bridge", "local_network_status", "mqtt_bridge", "firmware_ota"]' ;;
    sensor) echo '["temperature", "humidity", "motion", "battery", "firmware_ota"]' ;;
    switch_panel) echo '["multi_gang_switch", "scene_trigger", "led_indicator", "firmware_ota"]' ;;
    curtain) echo '["open", "close", "position", "calibration", "firmware_ota"]' ;;
    linux_simulator) echo '["virtual_device", "mqtt_payload_inspection", "debug_report", "firmware_ota"]' ;;
  esac
}

type_group_for_type() {
  case "$1" in
    light_bulb) echo "lighting" ;;
    air_conditioner) echo "hvac" ;;
    smart_plug) echo "energy" ;;
    camera) echo "video" ;;
    gateway) echo "gateway" ;;
    sensor) echo "sensing" ;;
    switch_panel) echo "scene-control" ;;
    curtain) echo "covering" ;;
    linux_simulator) echo "simulation" ;;
  esac
}

type_counts_file="$(mktemp)"
trap 'rm -f "$type_counts_file"' EXIT

for i in $(seq 1 "$device_count"); do
  device_id="$(printf 'rtk-test-device-%04d' "$i")"
  type="$(device_type_for_index "$i")"
  current_count="$(awk -F= -v t="$type" '$1 == t {print $2}' "$type_counts_file" 2>/dev/null | tail -1)"
  current_count="${current_count:-0}"
  ordinal=$((current_count + 1))
  grep -v "^$type=" "$type_counts_file" > "$type_counts_file.tmp" 2>/dev/null || true
  mv "$type_counts_file.tmp" "$type_counts_file"
  printf '%s=%s\n' "$type" "$ordinal" >> "$type_counts_file"

  model="$(model_for_type "$type")"
  display_name="$(display_name_for_type "$type" "$(printf '%02d' "$ordinal")")"
  group="$(type_group_for_type "$type")"
  capabilities="$(capabilities_for_type "$type")"
  device_dir="$out_dir/devices/$type/$device_id"
  bundle_dir="$out_dir/bundles/$type"
  mkdir -p "$device_dir" "$bundle_dir"

  openssl ecparam -name prime256v1 -genkey -noout -out "$device_dir/device.key"
  chmod 0600 "$device_dir/device.key"

  cat > "$device_dir/csr.conf" <<EOF
[ req ]
prompt = no
distinguished_name = dn
req_extensions = req_ext

[ dn ]
C = TW
O = Realtek Connect Plus Development
OU = Test Device
CN = $device_id

[ req_ext ]
subjectAltName = @alt_names

[ alt_names ]
DNS.1 = $device_id.test.realtek-connect.local
EOF

  openssl req -new -sha256 \
    -key "$device_dir/device.key" \
    -out "$device_dir/device.csr" \
    -config "$device_dir/csr.conf"

  cat > "$device_dir/cert-ext.conf" <<EOF
basicConstraints = CA:FALSE
keyUsage = critical, digitalSignature
extendedKeyUsage = clientAuth
subjectAltName = DNS:$device_id.test.realtek-connect.local
EOF

  openssl x509 -req -sha256 \
    -in "$device_dir/device.csr" \
    -CA "$out_dir/ca/dev-device-ca.crt" \
    -CAkey "$out_dir/ca/dev-device-ca.key" \
    -CAserial "$serial_file" \
    -CAcreateserial \
    -days "$device_valid_days" \
    -out "$device_dir/device.crt" \
    -extfile "$device_dir/cert-ext.conf"

  cat "$device_dir/device.crt" "$device_dir/device.key" > "$bundle_dir/$device_id.pem"
  chmod 0600 "$bundle_dir/$device_id.pem"

  cat > "$device_dir/metadata.json" <<EOF
{
  "device_id": "$device_id",
  "display_name": "$display_name",
  "product_type": "$type",
  "model": "$model",
  "firmware_version": "0.9.0-test",
  "capabilities": $capabilities,
  "test_group": "$group",
  "certificate_profile": "dev-device-mtls-client",
  "subject_uri": "urn:realtek-connect:test-device:$device_id",
  "production": false,
  "warning": "Test-only credential fixture. Do not use for production or staging device identity."
}
EOF

  printf '%s,%s,%s,"%s",%s,%s,%s\n' \
    "$device_id" "$type" "$model" "$display_name" \
    "devices/$type/$device_id/device.crt" \
    "devices/$type/$device_id/device.key" \
    "bundles/$type/$device_id.pem" >> "$manifest"

  if [[ "$json_array_started" == "true" ]]; then
    printf ',\n' >> "$json_manifest"
  fi
  json_array_started=true
  cat >> "$json_manifest" <<EOF
  {
    "device_id": "$device_id",
    "product_type": "$type",
    "model": "$model",
    "display_name": "$display_name",
    "certificate_path": "devices/$type/$device_id/device.crt",
    "key_path": "devices/$type/$device_id/device.key",
    "bundle_path": "bundles/$type/$device_id.pem"
  }
EOF
done

printf '\n]\n' >> "$json_manifest"

cat > "$out_dir/README.md" <<EOF
# Realtek Connect+ Test Device Credentials

This directory contains generated test-only device identities for development
and evaluation of Realtek Connect+ device onboarding, mTLS, fleet, OTA, and app
flows.

These credentials are intentionally not production credentials.

- Device count: $device_count
- Device key type: EC P-256
- Certificate profile: mTLS client certificate, \`extendedKeyUsage=clientAuth\`
- CA subject: \`Realtek Connect Plus Dev Device CA\`
- Device validity: $device_valid_days days
- Production use: forbidden

The private keys in this folder are public test fixtures once shared. Never
register these identities against production services or customer environments.

## Layout

- \`ca/dev-device-ca.crt\`: development CA certificate.
- \`ca/dev-device-ca.key\`: development CA private key for regenerating test fixtures.
- \`devices/<product_type>/<device_id>/device.key\`: per-device private key.
- \`devices/<product_type>/<device_id>/device.crt\`: per-device client certificate.
- \`devices/<product_type>/<device_id>/device.csr\`: per-device CSR.
- \`devices/<product_type>/<device_id>/metadata.json\`: device test metadata.
- \`bundles/<product_type>/<device_id>.pem\`: certificate + private key bundle.
- \`manifest.csv\` and \`manifest.json\`: generated inventory.

## Device Type Allocation

- 0001-0015: light bulbs
- 0016-0030: air conditioners
- 0031-0042: smart plugs
- 0043-0062: cameras / PRO2 demos
- 0063-0072: gateways
- 0073-0084: sensors
- 0085-0090: switch panels
- 0091-0096: curtains
- 0097-0100: Linux simulators
EOF

echo "generated $device_count test device credentials under $out_dir"
