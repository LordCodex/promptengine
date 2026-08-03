---
document_id: core-internationalization-and-localization
title: Internationalization (i18n) and Localization (l10n) Standard
ecosystem: cross-cutting
dependencies:
  - core-frontend-architecture
  - core-frontend-security
  - core-seo-engineering
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Internationalization (i18n) and Localization (l10n) Standard

## Purpose & Inheritance
This document defines core standards for building globally-ready applications. It inherits from the [Frontend Architecture Standard](23-frontend-architecture-standard.md) and the [SEO Engineering Standard](27-seo-engineering-standard.md), ensuring that codebases, interfaces, and data models cleanly support multiple languages, regional conventions, currencies, time zones, and cultures without core design alterations.

---

## 1. i18n First: Separation of Concerns

All applications must separate **executable logic** from **user-visible content**.
- **No Hardcoded Strings**: Never hardcode user-facing labels, placeholder texts, validation error messages, or email templates inside visual components or controllers.
- **Translation Resource Files**: Place all strings in dedicated locale files (e.g. JSON translation sheets, message bundles) mapped to unique keys.
- **Translatable Logic**: Use translation helper functions (e.g., `t('auth.login.submit')`) to resolve visual strings dynamically based on the active locale.

---

## 2. Right-to-Left (RTL) Language Layouts

Ensure interfaces function correctly when displaying RTL scripts (Arabic, Hebrew, Persian):
- **Dynamic Directionality**: Enforce document dir attributes: `<html dir="rtl">` or `<html dir="ltr">` dynamically computed from the active locale.
- **CSS Logical Properties**: Avoid hardcoded physical layouts (left/right). Use logical CSS properties to ensure mirror behaviors trigger automatically:
  - Prefer `margin-inline-start` / `margin-inline-end` over `margin-left` / `margin-right`.
  - Prefer `padding-inline-start` / `padding-inline-end` over `padding-left` / `padding-right`.
  - Prefer `border-start-start-radius` over `border-top-left-radius`.
  - Prefer `inset-inline-start` over `left`.

---

## 3. Locale-Aware Formatting

### 3a. Dates & Times
- **UTC Server Storage**: Store all timestamps, logs, and database records in UTC. Convert to the user's local time zone only at the UI rendering boundary.
- **Browser Locale API**: Never manually concatenate date blocks. Use built-in Intl APIs (e.g. `Intl.DateTimeFormat`) or locale-aware libraries:
  - Format 12-hour vs 24-hour clocks based on local preference.
  - Respect timezone parameters for bookings, reports, and calendar listings.

### 3b. Numbers & Metrics
- Format thousands and decimal separators dynamically (`Intl.NumberFormat`):
  - English: `1,234,567.89`
  - German: `1.234.567,89`
  - French: `1 234 567,89`
- **Units of Measurement**: Provide configuration fallbacks to toggle measurements between metric (kilometers, Celsius, kilograms) and imperial (miles, Fahrenheit, pounds) structures.

---

## 4. Currency and Multi-Currency Rules

### 4a. Currency Representation
- Format money values using standard currency rules (`Intl.NumberFormat` with `style: 'currency'`).
- Display explicit currency symbols or ISO codes where ambiguity exists. **Never assume `$` always represents US Dollars (USD)** — distinguish it from CAD, AUD, or SGD.

### 4b. Multi-Currency Operations
To prevent data corruption and accounting issues, enforce these architectural rules:
- **Financial Record Integrity**: Currency conversions for display purposes must never modify original stored ledger transactions.
- **Currency Roles**: Keep a clear separation between:
  - *Book currency*: The base ledger currency of the business account.
  - *Customer payment currency*: The currency charged to the buyer.
  - *Author/Vendor settlement currency*: The payout currency.
  - *Reporting currency*: The currency used in administrative analytics.
- **Exchange Rates**:
  - Converted amounts are displayed to the user as "estimates" unless the rate has been explicitly locked for the active transaction.
  - Explicitly record and save the exact exchange rate used, the timestamp, and source provider for every completed transaction.

---

## 5. Regional Adapters: Addresses, Phones, and Payments

### 5a. Address Configurations
Do not assume every country shares a standard `Street / City / State / ZIP` structure:
- Make fields like State and Province optional or dynamic.
- Allow variable postal code formats (not all countries use 5-digit zip codes).
- Use country-specific address templates or standard APIs (e.g. Google Address Autocomplete) to adapt input fields.

### 5b. Phone Numbers
- Force international dialing formatting (+ [Country Code] [Number]).
- Validate inputs dynamically using regional metadata (e.g. `libphonenumber`). Avoid enforcing fixed character length checks.

### 5c. Localized Payment Gateways
Present only the payment methods supported by the active user's region:
- **Nigeria**: *Paystack*, *Bank Transfer* options.
- **Europe**: *SEPA* direct debit, *Mister Cash*, *Giropay*.
- **Global**: *Stripe*, *Apple Pay*, *Google Pay*.

---

## 6. Text, Search, and Typography

- **UTF-8 Support**: Enforce UTF-8 encoding across all systems, database tables, HTTP headers, API payloads, and file uploads.
- **Accented Characters & Sorting**: Ensure database and frontend search queries handle Unicode normalization (comparing `é` and `e` matching rules), case-insensitivity, and locale-specific sorting parameters.
- **Multilingual Fonts**: Choose typefaces that support diverse script characters (Latin, Arabic, CJK - Chinese, Japanese, Korean) without falling back to broken glyph placeholders.

---

## 7. Global SEO and Accessibility

- **Crawlability Hooks**: Use `hreflang` tags on indexable pages (`<link rel="alternate" hreflang="es" href="https://es.example.com/">`) to point search engines to localized variants.
- **Localized Metadata**: Serve corresponding localized titles, meta descriptions, and Open Graph tags.
- **Accessibility**: Translated elements must preserve ARIA tags, labels, and reading direction behaviors for screen readers.

---

## 8. Framework Implementations

- **Vue / Nuxt**: Route and translate using `@nuxtjs/i18n`. Keep templates clean of raw strings by placing text properties under `messages/` JSON sheets.
- **React / Next.js**: Utilize Next.js localized subpaths routing. Keep component functions clean by loading translations using hooks (`useTranslation()`).
- **Flutter**: Configure localization delegates (`GlobalMaterialLocalizations`) and reference dynamically resolved AppLocalizations files.

---

## 9. Required AI Agent Guidelines

When building or updating features, AI agents must abide by the following default assumptions:
1. **Never assume USD** as the application's default currency.
2. **Never assume English** as the fallback language.
3. **Never assume `MM/DD/YYYY`** as the default date format.
4. **Never assume a single time zone** (e.g. server timezone).
5. **Always verify locale requirements** when designing inputs, forms, payments, or scheduling features.

---

## Review Checklist

Verify the application against this checklist for global readiness:
- [ ] **No hardcoded user-visible strings**: All text strings are extracted to translation sheets or bundles.
- [ ] **Locale-aware dates**: Dates, times, and calendar values use Intl APIs or locale-aware helpers.
- [ ] **Locale-aware numbers**: Percentage and decimal separators adapt correctly to active locales.
- [ ] **Locale-aware currency formatting**: Currency is formatted with correct symbol placement, precision, and currency codes.
- [ ] **Time zones handled correctly**: Server logic persists timestamps in UTC; client handles user timezone conversions.
- [ ] **UTF-8 supported throughout**: Databases, APIs, and file uploads enforce UTF-8 configurations.
- [ ] **RTL considered**: Screen layouts mirror correctly; stylesheets use CSS logical properties instead of left/right attributes.
- [ ] **Payment methods localized**: Region-appropriate payment gateways are displayed.
- [ ] **Address formats adaptable**: Address inputs are flexible and accommodate international formats.
- [ ] **Multilingual layouts supported**: Layouts, buttons, and labels expand or wrap gracefully when translation lengths change.
- [ ] **Exchange rates treated correctly**: Financial records preserve original base ledger currency and capture locked transaction rates.
- [ ] **Accessibility preserved**: Screen readers announce localized attributes and follow correct reading directions.

---

## References
- Frontend Architecture: [core/23-frontend-architecture-standard.md](23-frontend-architecture-standard.md)
- SEO Engineering: [core/27-seo-engineering-standard.md](27-seo-engineering-standard.md)
- Intl API Reference: [https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Intl](https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/Intl)
- W3C Internationalization: [https://www.w3.org/International/](https://www.w3.org/International/)
