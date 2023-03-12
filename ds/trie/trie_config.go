package trie

type trieConfig struct {
	// for boolean configurations
	flags int16

	// for value configurations
	// TODO: set default value
	keyLimit uint16
}

// var configFlagAllowEmpty int16 = (1 << 0)

func (c trieConfig) setFlags(ctx *TrieContext) {
	// Ex:
	// if ctx.AllowEmpty == true {
	// 	c.flags |= configFlagAllowEmpty
	// }
}

func CreateNewTrieConfig(ctx *TrieContext) trieConfig {
	config := trieConfig{keyLimit: ctx.KeyLimit}

	config.setFlags(ctx)
	return config
}

func (c trieConfig) isExceedingKeyLimit(key string) bool {
	len := len(key)
	return len > int(c.keyLimit)
}
