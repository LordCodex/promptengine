# PHP Laravel Sandbox Reference

This directory serves as the executable reference implementation of the PHP/Laravel playbook standards. It provides a functional, linted boilerplate mapping exactly to the conventions defined in:
- [php-conventions.md](../../stacks/php-laravel/php-conventions.md)
- [laravel-logic.md](../../stacks/php-laravel/laravel-logic.md)
- [laravel-routing.md](../../stacks/php-laravel/laravel-routing.md)

---

## Structure
- `/app/Actions/`: Single-purpose business transaction classes (Service actions).
- `/app/Http/Controllers/`: Slim controllers handling HTTP requests and routing bindings.
- `/app/Http/Requests/`: FormRequest parameter validators.

---

## Validation and Analysis
To run static analysis checks against this sandbox block:
```bash
# Set up vendor packages
composer install

# Execute PHPStan analyzer
vendor/bin/phpstan analyse app/
```
