namespace Infra.Pulumi;

public class StressConfig
{
    // Settings without defaults
    public string Regions { get; set; }
    public string PublicKey { get; set; }
    public string Environment { get; set; }    
    public string S3ClientBucketName { get; set; }
    public string S3LogBucketName { get; set; }
    public string CurrentRegion { get; set; }
    public List<string> AllowedCidrBlocks { get; set; }
    public int DesiredCapacity { get; set; }


    // Settings with defaults
    public string VpcCidrBlock { get; set; } = "10.10.0.0/16";
    public List<string> PrometheusAllowedCidrBlocks { get; set; } = new List<string> {"0.0.0.0/0"};
    public string AmiName { get; set; } = "stl-stress_test_loader-cluster";
    public string SourceAmiRegion { get; set; } = "us-west-2";
    public List<string> EgressAllowedCidrBlocks { get; set; } = new List<string> {"0.0.0.0/0"};
    public string InstanceType { get; set; } = "t4g.nano";
    public string TelegrafUsername { get; set; } = "admin";
    public string TelegrafPassword { get; set; } = "";
    public string TelegrafUrl { get; set; } = "";

    public int MinSize { get; set; } = 1;
    public int MaxSize { get; set; } = 10;
    public int UpScalingAdjustment { get; set; } = -1;
    public int DownScalingAdjustment { get; set; } = -1;
    public int AzCount { get; set; } = 2;
    public int StlPort { get; set; } = 9005;

    public bool PublicIpOnLaunch { get; set; } = true;
    public bool CreateS3Buckets { get; set; } = false;

    public StressConfig()
    {
    }
}