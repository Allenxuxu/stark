package registry

// Filter is used to filter a service during the selection process
type Filter func([]*Service) []*Service

// FilterEndpoint is an endpoint based Next Filter which will
// only return services with the endpoint specified.
func FilterEndpoint(name string) Filter {
	return func(old []*Service) []*Service {
		var services []*Service

		for _, service := range old {
			for _, ep := range service.Endpoints {
				if ep.Name == name {
					services = append(services, service)
					break
				}
			}
		}

		return services
	}
}

// FilterLabel is a label based Next Filter which will
// only return services with the label specified.
func FilterLabel(key, val string) Filter {
	return func(old []*Service) []*Service {
		var services []*Service

		for _, service := range old {
			serv := new(Service)
			var nodes []*Node

			for _, node := range service.Nodes {
				if node.Metadata == nil {
					continue
				}

				if node.Metadata[key] == val {
					nodes = append(nodes, node)
				}
			}

			// only add service if there's some nodes
			if len(nodes) > 0 {
				// copy
				*serv = *service
				serv.Nodes = nodes
				services = append(services, serv)
			}
		}

		return services
	}
}

// FilterVersion is a version based Next Filter which will
// only return services with the version specified.
func FilterVersion(version string) Filter {
	return func(old []*Service) []*Service {
		var services []*Service

		for _, service := range old {
			if service.Version == version {
				services = append(services, service)
			}
		}

		return services
	}
}
