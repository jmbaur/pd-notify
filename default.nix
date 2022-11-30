{ buildGoModule, ... }:
buildGoModule {
  pname = "pd-notify";
  version = "0.1.1";
  src = ./.;
  vendorSha256 = "sha256-kEeS+X45Mmo7yNrA0MpChhYguk+c3ayVSNO7i6/daJY=";
}
