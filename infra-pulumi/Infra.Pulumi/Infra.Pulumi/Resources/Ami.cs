using Pulumi;
using Aws = Pulumi.Aws;

namespace Infra.Pulumi.Resources;

class Ami : ComponentResource
{
    [Output]
    public Output<string> AmiId { get; set; }
    
    public Ami(string name, ComponentResourceOptions opts = null) : base("stl:aws:Ami", name, opts)
    {
        // Set up Config
        var config = new Config();

        // Set up AMI
        var stl = new Aws.Ec2.AmiCopy("stl", new Aws.Ec2.AmiCopyArgs
        {
            Name = config.Require("name"),
            SourceAmiId = config.Require("source_ami_id"),
            SourceAmiRegion = config.Require("source_ami_region")
        });

        this.AmiId = stl.Id;

        // Signal to the UI that this resource has completed construction.
        this.RegisterOutputs();
    }
}

