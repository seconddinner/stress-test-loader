using Pulumi;
using Aws = Pulumi.Aws;
using Pulumi.Aws.Ec2;

namespace Infra.Pulumi.Resources;

class S3 : ComponentResource
{
    public S3(StressConfig cfg, ComponentResourceOptions opts = null) : base("stl:aws:S3", "stl-s3", opts)
    {

        foreach (var bucketName in new List<string> { cfg.S3ClientBucketName, cfg.S3LogBucketName })
        {
            // https://github.com/pulumi/pulumi/issues/3388
            // Currently, Pulumi does not support "create a resource if it does not exist"
            // or "check if a resource exists and then ..."
            // So, we have to use a boolean to manually determine if we should create the s3 buckets

            if (cfg.CreateS3Buckets)
            {
                var s3Bucket = new Aws.S3.Bucket(bucketName, new Aws.S3.BucketArgs
                {
                    BucketName = bucketName,
                    Acl = "private",
                }, new CustomResourceOptions
                {
                    Parent = this,
                    RetainOnDelete = true
                });
            }
        }
        
        RegisterOutputs();
    }
}

