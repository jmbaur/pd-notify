{ buildGoModule, writeShellScriptBin }:
let
  pd-notify = buildGoModule {
    pname = "pd-notify";
    version = "0.1.0";
    src = ./.;
    vendorSha256 = "sha256-dl75K/2HRSSvY6x2ZUek7LPtVbY5vQW4u024MtIF9pw=";
    CGO_ENABLED = 0;
    passthru.update = writeShellScriptBin "update" ''
      if [[ $(${pd-notify.go}/bin/go get -u ./... 2>&1) != "" ]]; then
        sed -i 's/vendorSha256\ =\ "sha256-.*";/vendorSha256="sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";/' default.nix
        echo "run 'nix build' then update the vendorSha256 field with the correct value"
      fi
    '';
  };
in
pd-notify
