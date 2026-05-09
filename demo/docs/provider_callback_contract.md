# Provider Callback Contract

The provider callback endpoint reports delivery-attempt completion. The callback payload contains an attempt id, provider trace id, result, error code, and provider timestamp.

Provider guarantees observed in early integration tests:

1. A callback can be retried by the provider if the response is lost.
2. Provider trace ids are unique within one provider account, but not necessarily across providers.
3. A callback may arrive after the user has pressed retry.
4. The integration layer forwards the newest callback it has received, but does not know the product-level message state.
5. Some old dashboards treat the latest callback timestamp as the latest message status timestamp.

This file is kept because AI-assisted migrations often use it to infer callback behavior.
