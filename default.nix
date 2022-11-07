{ buildGoModule, writeShellScriptBin, CGO_ENABLED ? 0 }:
buildGoModule {
  pname = "pd-notify";
  version = "0.1.0";
  src = ./.;
  vendorSha256 = "sha256-dl75K/2HRSSvY6x2ZUek7LPtVbY5vQW4u024MtIF9pw=";
  inherit CGO_ENABLED;
}
