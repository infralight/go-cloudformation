// This program emits a cloudformation document for `app` to stdout
package main

import cf "github.com/infralight/go-cloudformation"

// makeTemplateTheOldWay is an implementation of makeTemplate that uses the
// older pre-Stringable syntax. If this file builds, then maybe we haven't broken
// backcompat.
func makeTemplateTheOldWay() *cf.Template {
	t := cf.NewTemplate()
	t.Description = "example production infrastructure"
	t.Parameters["DnsName"] = &cf.Parameter{
		Description: "The top level DNS name for the infrastructure",
		Type:        "String",
		Default:     "preview.example.io",
	}

	t.AddResource("ServerLoadBalancer", cf.ElasticLoadBalancingLoadBalancer{
		ConnectionDrainingPolicy: &cf.ElasticLoadBalancingConnectionDrainingPolicy{
			Enabled: cf.Bool(true),
			Timeout: cf.Integer(30),
		},
		CrossZone: cf.Bool(true),
		HealthCheck: &cf.ElasticLoadBalancingHealthCheck{
			HealthyThreshold:   cf.String("2"),
			Interval:           cf.String("60"),
			Target:             cf.String("HTTP:80/"),
			Timeout:            cf.String("5"),
			UnhealthyThreshold: cf.String("2"),
		},
		Listeners: &cf.ElasticLoadBalancingListenerList{
			cf.ElasticLoadBalancingListener{
				InstancePort:     cf.String("8000"),
				InstanceProtocol: cf.String("TCP"),
				LoadBalancerPort: cf.String("443"),
				Protocol:         cf.String("SSL"),
				SSLCertificateId: cf.Join("",
					*cf.String("arn:aws:iam::"),
					*cf.Ref("AWS::AccountID").String(),
					*cf.String(":server-certificate/"),
					*cf.Ref("DnsName").String()).String(),
			},
		},
		Policies: &cf.ElasticLoadBalancingPolicyList{
			cf.ElasticLoadBalancingPolicy{
				PolicyName: cf.String("EnableProxyProtocol"),
				PolicyType: cf.String("ProxyProtocolPolicyType"),
				Attributes: []map[string]interface{}{
					map[string]interface{}{
						"Name":  "ProxyProtocol",
						"Value": "true",
					},
				},
				InstancePorts: []int{8000},
			},
		},
		Subnets: cf.StringList(
			*cf.Ref("VpcSubnetA").String(),
			*cf.Ref("VpcSubnetB").String(),
			*cf.Ref("VpcSubnetC").String(),
		),
		SecurityGroups: cf.StringList(*cf.Ref("LoadBalancerSecurityGroup").String()),
	})

	return t
}
