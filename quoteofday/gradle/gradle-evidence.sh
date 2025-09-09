jf evd create --subject-repo-path evidence-demo-libs-release-local/com/example/quote-of-day-service/1.0.0-10/quote-of-day-service-1.0.0-10.jar \
  --key ../../../../evidence/CI-RSA-KEY.key --key-alias CI_RSA_KEY \
  --predicate gradle-build-tool.json --markdown gradle-build-tool.md \
  --predicate-type https://gradle.com/evidence/build-tool/v1 --provider-id gradle

jf evd create --subject-repo-path evidence-demo-libs-release-local/com/example/quote-of-day-service/1.0.0-10/quote-of-day-service-1.0.0-10.jar \
  --key ../../../../evidence/CI-RSA-KEY.key --key-alias CI_RSA_KEY \
  --predicate gradle-java-toolchain.json --markdown gradle-java-toolchain.md \
  --predicate-type https://gradle.com/evidence/java-toolchain/v1 --provider-id gradle 

jf evd create --subject-repo-path evidence-demo-libs-release-local/com/example/quote-of-day-service/1.0.0-10/quote-of-day-service-1.0.0-10.jar \
  --key ../../../../evidence/CI-RSA-KEY.key --key-alias CI_RSA_KEY \
  --predicate gradle-resolved-dependencies.json --markdown gradle-resolved-dependencies.md \
  --predicate-type https://gradle.com/evidence/resolved-dependencies/v1 --provider-id gradle