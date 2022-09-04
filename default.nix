{ buildGoModule, writeShellScriptBin, lib }:
let
  pd-notify = buildGoModule {
    pname = "pd-notify";
    version = "0.1.0";
    src = ./.;
    vendorSha256 = "sha256-twupKh4kIieHI9Zyg3OVzpy4f9dmMqPEhJ/red2e0Xk=";
    CGO_ENABLED = 0;
    passthru.update = writeShellScriptBin "update" ''
      if [[ $(${pd-notify.go}/bin/go get -u ./...) != "" ]]; then
        sed -i 's/vendorSha256\ =\ "sha256-.*";/vendorSha256="sha256-AAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAAA=";/' default.nix
        echo "run 'nix build' then update the vendorSha256 field with the correct value"
      fi
    '';
  };
in
pd-notify
