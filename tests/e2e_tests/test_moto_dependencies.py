import pytest


@pytest.mark.moto
async def test_moto_server_works(aws_credentials, s3_server, s3_client):
    bucket_name = "test-buck"
    await s3_client.create_bucket(Bucket=bucket_name, CreateBucketConfiguration=dict(LocationConstraint="eu-central-1"))
    await s3_client.put_object(Bucket=bucket_name, Key="keys/key1.pem", Body="Hello woulda!")
    result = await s3_client.list_objects_v2(Bucket=bucket_name, Prefix="keys/")
    keys = []
    for item in result["Contents"]:
        if item["Key"].endswith(".pem"):
            contents = await s3_client.get_object(Bucket=bucket_name, Key=item["Key"])
            keys.append((item["LastModified"], contents))
    contents = await keys[0][1]["Body"].read()
    assert contents.decode() == "Hello woulda!"