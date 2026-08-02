---
document_id: bridges-laravel-inertia-vue-engineering-standard
title: Laravel + Inertia.js + Vue 3 Engineering Standard
ecosystem: bridge
target_versions:
  laravel: "^11.0"
  inertia: "^1.0"
  vue: "^3.4"
dependencies:
  - core-universal-coding-standards
  - core-architecture-and-simplicity
  - core-database-engineering-standard
  - core-api-engineering-standard
  - core-security-engineering-standard
  - core-testing-engineering-standard
  - stacks-php-conventions
  - stacks-laravel-engineering-standard
  - stacks-vue-ts-engineering-standard
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Laravel + Inertia.js + Vue 3 Engineering Standard

## Purpose & Inheritance
This document defines the core standards for applications built on the **Laravel + Inertia.js + Vue 3** stack. It inherits from and extends the [Universal Coding Standards](../core/05-universal-coding-standards.md), the [Architecture Standards](../core/02-architecture-and-simplicity.md), the [Database Engineering Standard](../core/06-database-engineering-standard.md), the [API Engineering Standard](../core/07-api-engineering-standard.md), the [Security Engineering Standard](../core/08-security-engineering-standard.md), the [Laravel Engineering Standard](../stacks/php-laravel/laravel-engineering-standard.md), and the [Vue 3 and TypeScript Engineering Standard](../stacks/js-ts-vue-nuxt/vue-ts-engineering-standard.md). It establishes implementation protocols for full-stack monolithic applications utilizing Inertia.js as a stateless protocol bridge.

---

## 1. Architecture Philosophy

Inertia.js is not a traditional client-server API (like REST or GraphQL); it is a **stateless routing and data transport protocol** that connects a server-side framework (Laravel) with a client-side library (Vue 3).

### Architectural Boundaries
- **Laravel (The Server Host)**: Holds all application authority, routing, authentication validation, permissions policies, database transactions, database migrations, jobs queues, and core business algorithms.
- **Vue 3 (The Client Presenter)**: Serves strictly as a rendering component, user interface event handler, and presenter of UI layout states.
- **No Client-Side Business Logic**: UI templates must not compute interest rates, determine permission permissions (independent of server values), or build raw database updates. The frontend simply displays data and captures user events.

---

## 2. Directory Layouts

We organize our full-stack codebase to enforce clear separation of concerns.

### Project Folder Organization

```text
app/                     # Laravel Core Domain
├── Actions/             # Single-responsibility business routines
├── Models/              # Eloquent schemas and scopes
├── Services/            # Third-party integrations
├── Policies/            # Domain authorization rules
└── Http/
    ├── Controllers/     # Thin routing coordinators
    ├── Middleware/      # Request filters (InertiaShareMiddleware)
    └── Requests/        # Form validation schemas

resources/js/            # Vue Client Domain
├── Components/          # Shared visual elements (Buttons, Inputs)
├── Layouts/             # Shared page frames (AdminLayout.vue)
├── Pages/               # Router target components (Invoices/Index.vue)
├── Composables/         # Isolated browser state fakes
└── Stores/              # Global UI state (Pinia)
```

---

## 3. Routing & Controller Standards

### Routing Protocols
- **Single Source of Routing**: All routes must be declared in Laravel (`routes/web.php`). Do not implement client-side routers (like Vue Router).
- **Route Model Binding**: Always use type-hinted route model binding in controllers to automatically resolve database records:
  ```php
  Route::get('/invoices/{invoice}', [InvoiceController::class, 'show']);
  ```
- **Ziggy Routing**: Use the `route()` helper on the client side (powered by Ziggy) to generate routes dynamically from Laravel definitions. Do not hardcode URL paths in Vue components.

### Controller Coordination
Controllers must remain thin coordinate layers. Their sole responsibility is parsing request inputs, calling an Action or Service, and returning an Inertia response.

```php
// Good: Slim Controller invoking Single-Action and returning response
namespace App\Http\Controllers;

use App\Actions\PayInvoiceAction;
use App\Http\Requests\PayInvoiceRequest;
use App\Models\Invoice;
use Inertia\Inertia;
use Inertia\Response;
use Illuminate\Http\RedirectResponse;

class InvoicePaymentController extends Controller
{
    public function show(Invoice $invoice): Response
    {
        $this->authorize('view', $invoice);

        return Inertia::render('Invoices/Show', [
            'invoice' => new InvoiceResource($invoice),
            'stripeKey' => config('services.stripe.key')
        ]);
    }

    public function store(PayInvoiceRequest $request, Invoice $invoice, PayInvoiceAction $payInvoice): RedirectResponse
    {
        $this->authorize('update', $invoice);

        $payInvoice->execute($invoice, $request->validated());

        return redirect()->route('invoices.show', $invoice)
            ->with('flash.success', 'Invoice paid successfully.');
    }
}
```

---

## 4. Data Sharing & Props Design

### Shared Data (`Inertia::share`)
- **Shared Data Limits**: Keep shared data minimal. Only register variables that are accessed globally across the layout shell (e.g., authenticated user details, unread notification counts, current flash success alerts).
- **Use Lazy Shares**: Defer heavy database queries in shared keys by wrapping them in closures to prevent them from executing on every request.

### Props Scoping (Props Design)
- **Never Pass Naked Models**: Do not pass Eloquent models directly to page props. This exposes database structures and internal fields (like password hashes, API tokens, internal IDs) to the browser.
- **Enforce API Resources**: Transform all page props using Eloquent Resources to define exact output schemas:
  ```php
  // Good: Output attributes explicitly declared
  'invoice' => [
      'id' => $this->hashid,
      'amount' => $this->amount_cents,
      'status' => $this->status
  ]
  ```

---

## 5. Forms, Validation & Error Handling

### Form Validation
- **Server Validation is the Source of Truth**: Enforce Form Requests (`Http/Requests`) to validate all inputs on the server. Never rely on frontend Javascript validation for security.
- **Inertia useForm Hook**: Manage client-side form submissions, loading indicators, and validation error messages using the Inertia form utility:

```vue
<!-- Good: Inertia Typed Form Handling -->
<script setup lang="ts">
import { useForm } from '@inertiajs/vue3';

interface Props {
  invoiceId: string;
}

const props = defineProps<Props>();

const form = useForm({
  amount_cents: 0,
  payment_method: 'credit_card'
});

function submitForm() {
  form.post(route('invoices.pay', props.invoiceId), {
    preserveScroll: true,
    onSuccess: () => form.reset()
  });
}
</script>

<template>
  <form @submit.prevent="submitForm">
    <input type="number" v-model="form.amount_cents" />
    <span v-if="form.errors.amount_cents" class="error">{{ form.errors.amount_cents }}</span>

    <button type="submit" :disabled="form.processing">
      {{ form.processing ? 'Processing...' : 'Pay Invoice' }}
    </button>
  </form>
</template>
```

---

## 6. Authentication, Authorization & Security

Inertia runs inside standard session and cookie authentication structures.

### Security Configurations
- **Authorization Enforcement**: Frontend visibility controls (like hiding a "Delete" button from unauthorized users) are visual enhancements, not security. You must enforce permissions checks on the server using Laravel Policies:
  ```php
  // Good: Backend checks permissions explicitly
  $this->authorize('delete', $invoice);
  ```
- **CSRF Protection**: Ensure that Inertia is configured to automatically parse the `XSRF-TOKEN` cookie returned by Laravel sessions and append it to request headers.
- **File Upload Security**: Always validate uploaded files on the server using strict mime type and magic byte checks. Do not trust file extensions passed from the client browser.

---

## 7. Database Operations & Query Safety

When returning data to Inertia pages:
- **Eager Load Relationships**: Always eager load relationships in the controller to prevent $N+1$ query issues:
  ```php
  // Good: Eager load relations to prevent database connection loops
  $invoices = Invoice::with('customer')->paginate(15);
  ```
- **Cursor Pagination**: Use cursor pagination (`cursorPaginate`) instead of offset pagination (`paginate`) for high-volume list screens to improve query performance.

---

## 8. SRE Operations: Queues & Events

In monolithic architectures, background jobs and events keep application execution fast.

### SRE Integration Rules
- **Queues for Asynchronous Tasks**: Offload slow processes (such as email delivery, PDF report generation, third-party payment transactions) to background queues using Laravel Jobs.
- **Jobs Retries**: Specify retry limits (`public $tries = 3`) and dead-letter queue fallbacks to prevent failed jobs from blocking workers.
- **Decoupled Listeners**: Use Laravel Events and Listeners to decouple system layers (e.g., dispatch an `InvoicePaid` event, and let separate listeners send notifications and update analytics stores).

---

## 9. Performance Optimization (Laravel + Inertia)

### Inertia Optimization Strategies
- **Inertia Partial Reloads**: Configure page requests to reload only the required properties when a user interacts with a page, rather than fetching the entire dataset.
- **Lazy Props**: Wrap heavy database queries in `Inertia::lazy` closures to load them on demand (e.g., loading an audit history tab only when the user clicks it):
  ```php
  return Inertia::render('Invoices/Show', [
      'history' => Inertia::lazy(fn() => $invoice->history()->get())
  ]);
  ```
- **Webpack / Vite Code Splitting**: Ensure Vite is configured to compile page components into separate code-split bundles to reduce initial page load times.

---

## 10. Testing Strategy

### Dual-Layer Testing Protocol
- **Laravel Feature Tests**: Validate route authorizations, validation schemas, and database transactions. Assert that controllers return the correct page components and props structures:
  ```php
  public function test_invoice_details_page_renders_with_authorized_props()
  {
      $user = User::factory()->create();
      $invoice = Invoice::factory()->create(['user_id' => $user->id]);

      $this->actingAs($user)
          ->get(route('invoices.show', $invoice))
          ->assertInertia(fn (Assert $page) => $page
              ->component('Invoices/Show')
              ->has('invoice')
              ->where('invoice.id', $invoice->hashid)
          );
  }
  ```
- **Vue Interaction Testing**: Test component behavioral interactions, events emission, and state modifications using Vitest and Vue Test Utils (see the [Vue standard](../stacks/js-ts-vue-nuxt/vue-ts-engineering-standard.md)).

---

## 11. Legacy Refactoring

When refactoring legacy Laravel systems containing mixed Blade files and large database-heavy controllers:
1. **Define the Target Boundary**: Convert Blade views to Inertia pages page-by-page. Avoid massive rewrites of the entire UI at once.
2. **Move Business Logic to Actions**: Extract complex logic from legacy controllers into single-responsibility Actions before rewriting views.
3. **Verify API Contract**: Ensure that the props returned by your new Inertia controllers exactly match the interfaces expected by the corresponding Vue page components.

---

## 12. Decision Matrices

Use these matrices to identify the correct bridge engineering decision based on project context.

### Matrix 1: Controller vs. Service vs. Action
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Coordinate HTTP requests, validate inputs, route responses | **Controller** | Standard thin routing coordinator. |
| Third-party API adapters (Stripe, Twilio) | **Service** | Centralizes SDK initializations and client connections. |
| Single business workflow (e.g., executing a checkout) | **Action** | Reusable, single-responsibility class that isolates business logic. |

### Matrix 2: Inertia Page Props vs. Async API Endpoint
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Initial data required to render a page layout | **Inertia Props** | Simplifies data flow by passing data alongside the page rendering. |
| Frequently updating lists, search autocomplete, live updates | **API Endpoint** | Fetches data asynchronously, avoiding full-page reloads. |

### Matrix 3: Shared Data vs. Page Props
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| User profile info, unread messages count, permissions | **Shared Data** | Globally accessible across layouts. |
| Specific invoice records, product lists, audit histories | **Page Props** | Keeps payloads scoped to relevant routing targets. |

### Matrix 4: Event vs. Direct Action Call
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Critical business operations (e.g. updating payment records) | **Direct Call** | Simple data flow; ensures immediate database consistency. |
| Secondary side-effects (e.g., sending emails, Slack alerts) | **Event / Listener** | Decouples system layers and prevents slow external API calls from blocking execution. |

### Matrix 5: Queue vs. Immediate Execution
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Direct database inserts, validating session credentials | **Immediate** | Critical path required for instant page feedback. |
| PDF compilation, exporting records lists, batch imports | **Queue Job** | Offloads slow operations to background workers, keeping the UI responsive. |

### Matrix 6: Vue Ref State vs. Inertia Form Server State
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Visual UI toggles, modal states, sidebar expands | **Vue Ref State** | Fast, client-side visual changes. |
| Text inputs, password resets, checkbox choices | **Inertia Form State** | Automatically binds inputs, tracks processing states, and maps server errors. |

---

## 13. AI Inertia Rules

AI agents modifying or writing code in this stack must follow these rules:

1. **Keep Controllers Slim**: Never write database queries or business calculations inside controllers. Delegate tasks to Actions or Eloquent scopes.
2. **Never Expose Naked Models**: Transform all Eloquent models into formatted resource arrays before passing them to Inertia render calls.
3. **No Direct DB Calls in Frontend**: Restrict all database transactions to the Laravel layer. The Vue client must call backend routes.
4. **Use Lazy Props for Deferral**: Wrap slow data requests in `Inertia::lazy` closures to keep initial page renders fast.
5. **Always Assert Inertia Components**: Ensure Laravel feature tests explicitly verify the target Inertia page component name and prop properties.

---

## 14. Monolithic Code Review Checklist

Use this checklist during code review to evaluate Laravel, Inertia, and Vue changes.

### Architecture & Boundaries
- [ ] Are controllers kept slim (delegating logic to Actions or Services)?
- [ ] Is frontend UI state kept separate from server business rules?

### Routing & Controllers
- [ ] Are routes declared exclusively in Laravel (`routes/web.php`)?
- [ ] Are route helpers generated dynamically using the Ziggy library (no hardcoded URLs)?
- [ ] Do controller methods authorize requests using Laravel Policies?

### Data Sharing & Page Props
- [ ] Are Eloquent models transformed using resources (no naked models passed as props)?
- [ ] Is shared data kept minimal and wrapped in lazy closures where appropriate?

### Forms & Validation
- [ ] Are form requests validated on the server using Laravel Form Requests?
- [ ] Do frontend forms handle loading states and error mapping using the `useForm` hook?

### Database & Performance
- [ ] Are Eloquent relations eager loaded (`with`) to prevent N+1 query loops?
- [ ] Have heavy database props been configured as `Inertia::lazy`?

### Testing
- [ ] Do feature tests verify target Inertia component names and props values?
- [ ] Do backend validation and policy rules have dedicated test coverage?

---

## References
- Secure Database Schemas: [core/06-database-engineering-standard.md](../core/06-database-engineering-standard.md)
- API Validation Envelope: [core/07-api-engineering-standard.md](../core/07-api-engineering-standard.md)
- Automated CI Pipelines: [core/13-cicd-and-deployment-standard.md](../core/13-cicd-and-deployment-standard.md)
- Vue Component Conventions: [stacks/js-ts-vue-nuxt/vue-ts-engineering-standard.md](../stacks/js-ts-vue-nuxt/vue-ts-engineering-standard.md)
