using Pulumi;
using Aws = Pulumi.Aws;
using Pulumi.Aws.Ec2;

namespace Infra.Pulumi.Resources;

class Ami : ComponentResource
{
    [Output]
    public Output<string> AmiId { get; set; }
    
    public Ami(string name, Aws.Provider provider, string region, ComponentResourceOptions opts = null) : base("stl:aws:Ami", name, opts)
    {
        // Set up Config
        var config = new Config();

        // Set up AMI
        var amiFilter = new Aws.Ec2.Inputs.GetAmiFilterArgs
        {
            Name = "name",
            Values = { "stress-test-loader*" }
        };

        var amiId = Output.Create(Aws.Ec2.GetAmi.InvokeAsync(new GetAmiArgs
        {
            Filters = { amiFilter },
            Owners = { "self" },
            MostRecent = true
        })).Apply(result => result.Id);

        var stl = new Aws.Ec2.AmiCopy("stl-" + region, new Aws.Ec2.AmiCopyArgs
        {
            Name = config.Require("name"),
            SourceAmiId = amiId,
            SourceAmiRegion = config.Require("source_ami_region")
        }, new CustomResourceOptions
        {
            Provider = provider,
            Parent = this
        });

        this.AmiId = stl.Id;
        
        RegisterOutputs();
    }
}

