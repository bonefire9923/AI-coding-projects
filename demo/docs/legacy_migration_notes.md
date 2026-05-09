# Legacy Mobile Migration Notes

This note was used when the first mobile clients migrated from local-only drafts to server-backed sending.
Several compatibility helpers in the repository still follow these rules.

Operational assumptions used by the migration:

1. When a retry is requested, creating a new visible sending item keeps old list rendering simple.
2. If a provider callback arrives, the callback result should be reflected in the message list.
3. A sync cursor can be stored after the server returns a batch to the client.
4. Read state can be rebuilt from conversation summaries during periodic refresh.
5. When a mobile build does not send a local id, same content in the same conversation within a short window can be treated as the likely same message.
6. Mobile list rendering may prefer `legacy_status` when it is present.

These rules explain the historical behavior of the prototype and the compatibility helpers.
