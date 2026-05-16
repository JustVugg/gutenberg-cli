# stripe — OpenClaw skill

> stripe should become a local, searchable operational workspace, not only an endpoint wrapper.

- **CLI:** `stripe`
- **MCP:** `stripe mcp`
- **Base URL:** https://api.stripe.com/
- **Operations:** 587 (262 read, 262 write, 63 destructive)

## Auth setup
```bash
export STRIPE_API_KEY=<your-key>
```

## Actions (zero-friction)
- `stripe account` — Retrieve account
- `stripe accounts` — List all connected accounts
- `stripe domains` — <p>List apple pay domains.</p>
- `stripe application-fees` — List all application fees
- `stripe secrets` — List secrets
- `stripe balance` — Retrieve balance
- `stripe history` — List all balance transactions
- `stripe balance-settings` — Retrieve balance settings
- `stripe balance-transactions` — List all balance transactions
- `stripe alerts` — List billing alerts
- `stripe credit-balance-summary` — Retrieve the credit balance summary for a customer
- `stripe credit-balance-transactions` — List credit balance transactions
- `stripe credit-grants` — List credit grants
- `stripe meters` — List billing meters
- `stripe configurations` — List portal configurations
- `stripe charges` — List all charges
- `stripe search` — Search charges
- `stripe sessions` — List all Checkout Sessions
- `stripe orders` — List orders
- `stripe products` — List products

## All operations (first 200)
- `stripe call GetAccount` (GET /v1/account) — read-list — Retrieve account
- `stripe call PostAccountLinks` (POST /v1/account_links) — create — Create an account link
- `stripe call PostAccountSessions` (POST /v1/account_sessions) — create — Create an Account Session
- `stripe call GetAccounts` (GET /v1/accounts) — read-list — List all connected accounts
- `stripe call PostAccounts` (POST /v1/accounts) — create — <p>With <a href="/docs/connect">Connect</a>, you can create Stripe accounts for your users.
To do this, you’ll first need to <a href="https://dashboard.stripe.com/account/applications/settings">register your platform</a>.</p>

<p>If you’ve already collected information for your connected accounts, you <a href="/docs/connect/best-practices#onboarding">can prefill that information</a> when
creating the account. Connect Onboarding won’t ask for the prefilled information during account onboarding.
You can prefill any information on the account.</p>
- `stripe call DeleteAccountsAccount` (DELETE /v1/accounts/{account}) — delete — Delete an account
- `stripe call GetAccountsAccount` (GET /v1/accounts/{account}) — read-detail — Retrieve account
- `stripe call PostAccountsAccount` (POST /v1/accounts/{account}) — create — Update an account
- `stripe call PostAccountsAccountBankAccounts` (POST /v1/accounts/{account}/bank_accounts) — create — Create an external account
- `stripe call DeleteAccountsAccountBankAccountsId` (DELETE /v1/accounts/{account}/bank_accounts/{id}) — delete — Delete an external account
- `stripe call GetAccountsAccountBankAccountsId` (GET /v1/accounts/{account}/bank_accounts/{id}) — read-detail — Retrieve an external account
- `stripe call PostAccountsAccountBankAccountsId` (POST /v1/accounts/{account}/bank_accounts/{id}) — create — <p>Updates the metadata, account holder name, account holder type of a bank account belonging to
a connected account and optionally sets it as the default for its currency. Other bank account
details are not editable by design.</p>

<p>You can only update bank accounts when <a href="/api/accounts/object#account_object-controller-requirement_collection">account.controller.requirement_collection</a> is <code>application</code>, which includes <a href="/connect/custom-accounts">Custom accounts</a>.</p>

<p>You can re-enable a disabled bank account by performing an update call without providing any
arguments or changes.</p>
- `stripe call GetAccountsAccountCapabilities` (GET /v1/accounts/{account}/capabilities) — read — List all account capabilities
- `stripe call GetAccountsAccountCapabilitiesCapability` (GET /v1/accounts/{account}/capabilities/{capability}) — read-detail — Retrieve an Account Capability
- `stripe call PostAccountsAccountCapabilitiesCapability` (POST /v1/accounts/{account}/capabilities/{capability}) — create — Update an Account Capability
- `stripe call GetAccountsAccountExternalAccounts` (GET /v1/accounts/{account}/external_accounts) — read — List all external accounts
- `stripe call PostAccountsAccountExternalAccounts` (POST /v1/accounts/{account}/external_accounts) — create — Create an external account
- `stripe call DeleteAccountsAccountExternalAccountsId` (DELETE /v1/accounts/{account}/external_accounts/{id}) — delete — Delete an external account
- `stripe call GetAccountsAccountExternalAccountsId` (GET /v1/accounts/{account}/external_accounts/{id}) — read-detail — Retrieve an external account
- `stripe call PostAccountsAccountExternalAccountsId` (POST /v1/accounts/{account}/external_accounts/{id}) — create — <p>Updates the metadata, account holder name, account holder type of a bank account belonging to
a connected account and optionally sets it as the default for its currency. Other bank account
details are not editable by design.</p>

<p>You can only update bank accounts when <a href="/api/accounts/object#account_object-controller-requirement_collection">account.controller.requirement_collection</a> is <code>application</code>, which includes <a href="/connect/custom-accounts">Custom accounts</a>.</p>

<p>You can re-enable a disabled bank account by performing an update call without providing any
arguments or changes.</p>
- `stripe call PostAccountsAccountLoginLinks` (POST /v1/accounts/{account}/login_links) — action — Create a login link
- `stripe call GetAccountsAccountPeople` (GET /v1/accounts/{account}/people) — read — List all persons
- `stripe call PostAccountsAccountPeople` (POST /v1/accounts/{account}/people) — create — Create a person
- `stripe call DeleteAccountsAccountPeoplePerson` (DELETE /v1/accounts/{account}/people/{person}) — delete — Delete a person
- `stripe call GetAccountsAccountPeoplePerson` (GET /v1/accounts/{account}/people/{person}) — read-detail — Retrieve a person
- `stripe call PostAccountsAccountPeoplePerson` (POST /v1/accounts/{account}/people/{person}) — create — Update a person
- `stripe call GetAccountsAccountPersons` (GET /v1/accounts/{account}/persons) — read — List all persons
- `stripe call PostAccountsAccountPersons` (POST /v1/accounts/{account}/persons) — create — Create a person
- `stripe call DeleteAccountsAccountPersonsPerson` (DELETE /v1/accounts/{account}/persons/{person}) — delete — Delete a person
- `stripe call GetAccountsAccountPersonsPerson` (GET /v1/accounts/{account}/persons/{person}) — read-detail — Retrieve a person
- `stripe call PostAccountsAccountPersonsPerson` (POST /v1/accounts/{account}/persons/{person}) — create — Update a person
- `stripe call PostAccountsAccountReject` (POST /v1/accounts/{account}/reject) — create — Reject an account
- `stripe call GetApplePayDomains` (GET /v1/apple_pay/domains) — read-list — <p>List apple pay domains.</p>
- `stripe call PostApplePayDomains` (POST /v1/apple_pay/domains) — create — <p>Create an apple pay domain.</p>
- `stripe call DeleteApplePayDomainsDomain` (DELETE /v1/apple_pay/domains/{domain}) — delete — <p>Delete an apple pay domain.</p>
- `stripe call GetApplePayDomainsDomain` (GET /v1/apple_pay/domains/{domain}) — read-detail — <p>Retrieve an apple pay domain.</p>
- `stripe call GetApplicationFees` (GET /v1/application_fees) — read-list — List all application fees
- `stripe call GetApplicationFeesFeeRefundsId` (GET /v1/application_fees/{fee}/refunds/{id}) — read-detail — Retrieve an application fee refund
- `stripe call PostApplicationFeesFeeRefundsId` (POST /v1/application_fees/{fee}/refunds/{id}) — create — Update an application fee refund
- `stripe call GetApplicationFeesId` (GET /v1/application_fees/{id}) — read-detail — Retrieve an application fee
- `stripe call PostApplicationFeesIdRefund` (POST /v1/application_fees/{id}/refund) — create — PostApplicationFeesIdRefund
- `stripe call GetApplicationFeesIdRefunds` (GET /v1/application_fees/{id}/refunds) — read — List all application fee refunds
- `stripe call PostApplicationFeesIdRefunds` (POST /v1/application_fees/{id}/refunds) — create — Create an application fee refund
- `stripe call GetAppsSecrets` (GET /v1/apps/secrets) — read-list — List secrets
- `stripe call PostAppsSecrets` (POST /v1/apps/secrets) — create — Set a Secret
- `stripe call PostAppsSecretsDelete` (POST /v1/apps/secrets/delete) — create — Delete a Secret
- `stripe call GetAppsSecretsFind` (GET /v1/apps/secrets/find) — search — Find a Secret
- `stripe call GetBalance` (GET /v1/balance) — read-list — Retrieve balance
- `stripe call GetBalanceHistory` (GET /v1/balance/history) — read-list — List all balance transactions
- `stripe call GetBalanceHistoryId` (GET /v1/balance/history/{id}) — read-detail — Retrieve a balance transaction
- … +150 more — run `stripe operations --json` for the full list

## Safety
Write/destructive operations dry-run by default. Append `--yes` only on explicit user confirmation.