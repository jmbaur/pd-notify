{ buildGoModule, writeShellScriptBin, CGO_ENABLED ? 0 }:
buildGoModule {
  pname = "pd-notify";
  version = "0.1.0";
  src = ./.;
  vendorSha256 = "sha256-QJenPUIhfIXDHfCHoxX6CqOn+MGZrQPL97tvYP6iNno=";
  inherit CGO_ENABLED;
}
