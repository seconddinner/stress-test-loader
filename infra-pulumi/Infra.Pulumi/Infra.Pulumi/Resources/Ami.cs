using Pulumi;
using Aws = Pulumi.Aws;
using Pulumi.Aws.Ec2;
using Infra.Pulumi;

namespace Infra.Pulumi.Resources;

class Ami : ComponentResource
{
    public Output<string> AmiId { get; set; }
    
    public Ami(StressConfig cfg, ComponentResourceOptions opts = null) : base("stl:aws:Ami", $"stl-ami-{cfg.CurrentRegion}", opts)
    {

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

        var stl = new Aws.Ec2.AmiCopy("stl-" + cfg.CurrentRegion, new Aws.Ec2.AmiCopyArgs
        {
            Name = cfg.AmiName,
            SourceAmiId = amiId,
            SourceAmiRegion = cfg.SourceAmiRegion
        }, new CustomResourceOptions
        {
            Parent = this
        });

        this.AmiId = stl.Id;
        
        RegisterOutputs();
    }
}

