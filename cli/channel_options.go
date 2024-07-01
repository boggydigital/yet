package cli

type ChannelOptions struct {
	Force bool
}

func DefaultChannelOptions() *ChannelOptions {
	return &ChannelOptions{
		Force: false,
	}
}
