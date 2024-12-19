# GCloud Storage to S3 Migrator

## Configuration

### GCP

Use the default credentials for the GCP SDK: https://cloud.google.com/docs/authentication/provide-credentials-adc#how-to

Configure `GCLOUD_BUCKET_NAME` to specify the bucket to migrate from.

### S3

- `S3_ENDPOINT`: the endpoint for the s3 service
- `S3_SECRET_KEY`: the secret key for the s3 service
- `S3_SECRET_ACCESS_KEY`: the secret access key for the s3 service
- `S3_BUCKET_NAME`: the s3 bucket to migrate to
