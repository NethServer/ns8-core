module github.com/NethServer/ns8-scratchpad/core/agent

go 1.16

require github.com/go-redis/redis/v8 v8.8.0

replace github.com/NethServer/ns8-scratchpad/core/agent/models => ./models
replace github.com/NethServer/ns8-scratchpad/core/agent/action => ./action
