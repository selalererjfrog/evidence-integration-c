jf evd create --release-bundle quotopia-app --release-bundle-version 6 --project quotopia \
  --key ../../../../evidence/CI-RSA-KEY.key --key-alias CI_RSA_KEY \
  --predicate argocd.json --markdown argocd.md \
  --predicate-type https://jfrog.com/evidence/argocd/v1 --provider-id akuity