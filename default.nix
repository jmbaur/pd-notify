{ buildGoModule }:
buildGoModule {
  pname = "pd-notify";
  version = "0.1.0";
  src = ./.;
  vendorSha256 = "sha256-lvrleS3aNwQj+lpGUNhSi9+jfMueZdVz1ss+xtVGt+s=";
  CGO_ENABLED = 0;
}
