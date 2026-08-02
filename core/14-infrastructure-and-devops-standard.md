---
document_id: core-infrastructure-and-devops-standard
title: Infrastructure and DevOps Engineering Standard
ecosystem: cross-cutting
dependencies:
  - core-universal-coding-standards
  - core-architecture-and-simplicity
  - core-database-engineering-standard
  - core-api-engineering-standard
  - core-security-engineering-standard
  - core-security-testing-and-threat-modeling
  - core-performance-engineering-standard
  - core-testing-engineering-standard
  - core-git-and-collaboration-standard
  - core-cicd-and-deployment-standard
  - stacks-php-conventions
  - stacks-laravel-engineering-standard
audience: [human, agent]
last_reviewed: 2026-08-01
---

# Infrastructure and DevOps Engineering Standard

## Purpose & Inheritance
This document defines the core standards for system administration, Linux hardening, network configuration, process management, container orchestration, backup scheduling, observability, and disaster recovery. It inherits from and extends the [Universal Coding Standards](05-universal-coding-standards.md), the [Architecture Standards](02-architecture-and-simplicity.md), and all preceding core engineering documents. It establishes strict operational guidelines for human engineers and AI agents configuring infrastructure.

---

## 1. Infrastructure Philosophy

Infrastructure is an **extension of the application data layer**. Software cannot run reliably on poorly configured or un-monitored servers.
- **Simplicity Over Novelty**: We reject overengineering (e.g., using Kubernetes or distributed service meshes for simple CRUD apps). Avoid adding infrastructure layers unless active traffic volume demands scaling lanes.
- **Infrastructure as Declarative Code**: All environment setups, networking rules, and container configs must be version-controlled (Docker Compose, Terraform, Ansibles) rather than executed manually in bash terminals.
- **Secure by Default**: Block all entry ports by default. Allow incoming network access only through verified proxy/load-balancer boundaries.

---

## 2. Server & Linux Systems Administration

Linux (specifically Alpine and Ubuntu LTS distributions) is our standard production OS runtime environment.

### Linux Administration Directives
- **Low-Privilege Execution**: Never run application processes as `root`. Provision dedicated, non-login users (`appuser`) for runtime tasks.
- **Package Updates Schedule**: Enforce automatic installation of security updates (e.g., `unattended-upgrades` on Ubuntu).
- **Resource Limits**: Set system descriptor boundaries (`ulimit -n 65535`) to prevent sockets exhaustion during load spikes.
- **Monitoring Resource Usage**: Enforce system monitors to track metrics: Memory (OOM hazards), CPU usage, Disk space (logs exhaustion), and Input/Output operations (IOPS).

---

## 3. Linux Security & Hardening

Hardening Linux prevents unauthorized server access.

### Security Controls
1. **SSH Hardening**:
   - Disable password-based authentication (`PasswordAuthentication no` in `/etc/ssh/sshd_config`).
   - Force SSH key authentication using `RSA-4096` or `ED25519` key formats.
   - Disable root logins (`PermitRootLogin no`).
   - Change default SSH port (e.g., to a high port number like `2222`) to stop automated brute-force bots.
2. **Firewall Configurations**:
   - Enforce UFW (Uncomplicated Firewall) or Iptables.
   - Drop all incoming traffic by default. Open only required ports: SSH (`2222`), HTTP (`80`), and HTTPS (`443`).
3. **Sudo Isolation**: Restrict `sudo` access. Only allow key system administrators to call elevated commands, and require password checks for all sudo calls.
4. **Brute Force Protection**: Install `Fail2ban` to automatically block IP addresses that trigger repeated authentication failures on public SSH ports.
5. **Security Logging**: Ship system logs (`/var/log/auth.log`, `/var/log/syslog`) to a remote, write-once logging server for auditing.

---

## 4. Networking Foundations & Web Servers

Understanding network protocols is essential to resolve connection errors.

### Networking Essentials
- **TCP/UDP**: Use TCP for reliable, connection-oriented data (APIs, DBs, web traffic). Use UDP for speed-oriented, loss-tolerant streams (telemetry, DNS).
- **DNS Resolution**: Map records accurately:
  - `A`: Maps domain names to IPv4 addresses.
  - `AAAA`: Maps domain names to IPv6 addresses.
  - `CNAME`: Canonical name alias routing (pointing a subdomain to another domain).
  - `MX`: Mail exchange record (defines mail routing hosts).
  - `TXT`: Arbitrary text data (used for SPF, DKIM, and site verification).
- **Connectivity Troubleshooting**: Use diagnostic tools to isolate connection errors:
  - Check DNS resolution: `dig +trace domain.com`
  - Trace route paths: `traceroute -I destination_ip`
  - Probe ports: `nc -zv destination_ip port` or `curl -I https://destination_ip`

### Nginx Configuration Standard
Nginx is our default reverse proxy and SSL termination layer.

```text
HTTP Request (Port 443) ──> Nginx (SSL Termination & Rate Limiter)
                              └── Reverse Proxy (Port 9000) ──> PHP-FPM / Node App
```

#### Production Nginx Configuration Blueprint
```nginx
user nginx;
worker_processes auto;
error_log /var/log/nginx/error.log warn;
pid /var/run/nginx.pid;

events {
    worker_connections 2048;
}

http {
    include /etc/nginx/mime.types;
    default_type application/octet-stream;

    # SSL TLS Hardening
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers on;
    ssl_ciphers 'ECDHE-ECDSA-AES256-GCM-SHA384:ECDHE-RSA-AES256-GCM-SHA384';
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 1d;

    # Rate Limiting
    limit_req_zone $binary_remote_addr zone=api_limit:10m rate=30r/s;

    server {
        listen 80 default_server;
        server_name _;
        return 301 https://$host$request_uri; # Force HTTPS
    }

    server {
        listen 443 ssl http2;
        server_name api.domain.com;

        ssl_certificate /etc/letsencrypt/live/api.domain.com/fullchain.pem;
        ssl_certificate_key /etc/letsencrypt/live/api.domain.com/privkey.pem;

        # Security Headers
        add_header X-Frame-Options "DENY" always;
        add_header X-Content-Type-Options "nosniff" always;
        add_header Referrer-Policy "strict-origin-when-cross-origin" always;
        add_header Strict-Transport-Security "max-age=63072000; includeSubDomains; preload" always;

        location / {
            limit_req zone=api_limit burst=20 nodelay;
            
            proxy_pass http://127.0.0.1:9000;
            proxy_set_header Host $host;
            proxy_set_header X-Real-IP $remote_addr;
            proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
            proxy_set_header X-Forwarded-Proto $scheme;
        }
    }
}
```

- **SSL/TLS Certificates**: Secure all connections via HTTPS. Use **Let's Encrypt** with automated renewal cron jobs. Enforce TLS 1.2 or 1.3 only; disable deprecated TLS 1.0 and 1.1 protocol handshakes.
- **Apache Web Server**: Avoid using Apache unless legacy environments require it (e.g. relying on `.htaccess` configurations). Nginx offers significantly lower memory footprint and superior static file throughput.

---

## 5. Application Runtimes & Process Managers

Managed processes ensure that applications remain running after system reboots.

### Process Management Tools Matrix
- **Systemd**: Use Systemd to manage system-level daemons and global network services (e.g., Nginx, Redis, database engines).
- **Supervisor**: Best for queue workers and worker processes in PHP/Laravel environments (e.g., running `artisan queue:work` loops). It automatically restarts failed processes.
- **PM2**: Best for Node.js runtimes. Manages cluster modes, log rotations, and restart logic.

#### Supervisor Queue Worker Configuration Example
```ini
[program:laravel-worker]
process_name=%(program_name)s_%(process_num)02d
command=php /var/www/html/artisan queue:work --sleep=3 --tries=3 --max-time=3600
autostart=true
autorestart=true
stopasgroup=true
killasgroup=true
user=appuser
numprocs=4
redirect_stderr=true
stdout_logfile=/var/log/supervisor/laravel-worker.log
```

---

## 6. Database Operations SRE

Production database operations focus on data durability and performance safety.

### SRE Database Standards
- **Backup Schedules & Frequency**: Enforce automated, incremental daily backups and full weekly snapshots.
- **Storage Isolation**: Backups must be copied to isolated cloud storage buckets (e.g., AWS S3 with Object Lock enabled) in a separate geographical region. Never store database backups on the same physical disk as the active database.
- **Restore Drills**: Conduct monthly restoration testing drills to verify that backup files can be recovered and booted successfully.
- **Connection Tuning**: Enforce limits on active connection pools. Set up database connection proxy layers (e.g., PgBouncer for PostgreSQL) to prevent connection timeouts during application load spikes.

---

## 7. Container Orchestration

Container orchestration selection must align with actual traffic scale and operational budgets.

### Container Rules
- **Docker Compose**: Default option for local development, staging platforms, and single-host VPS production systems. It provides low operational overhead and reliable container networking.
- **Kubernetes (K8s)**: Enforce a strict **do not use Kubernetes by default** rule:
  - Do not use K8s unless your application handles high-concurrency requests (>10,000 requests/second), has dedicated SRE team support, and runs on multiple autoscaling cloud server groups.
  - Using Kubernetes on simple web apps adds unnecessary complexity, increases cloud costs, and degrades developer speed.
- **Managed Container Platforms**: Use managed container environments (e.g., AWS ECS, Google Cloud Run) to host microservice containers, rather than managing custom Kubernetes clusters on VM nodes.

---

## 8. Scaling, Backups & Observability

### Server Scaling Boundaries
- **Vertical Scaling**: Upgrade VM resources (CPU, RAM, disk speed) first. Best for database servers.
- **Horizontal Scaling**: Add more application server instances behind a load balancer. Requires that the application is completely stateless:
  - Sessions must reside in shared stores (e.g., Redis).
  - Uploaded files must be saved to cloud object storage (e.g., S3).

### Disaster Recovery Targets
- **Recovery Point Objective (RPO)**: The maximum age of data that can be lost during an outage (e.g., RPO = 24 hours means daily backups).
- **Recovery Time Objective (RTO)**: The maximum duration of system downtime allowed before services are recovered (e.g., RTO = 2 hours).

### Observability & Incident Response
- **Logging**: Collect system and application logs in a centralized location (e.g., ELK stack, Grafana Loki).
- **Alerting Rules**: Enforce alerts for critical events: disk usage $>80\%$, CPU load $>90\%$ for 5 consecutive minutes, API error rates $>1\%$.
- **Incident Post-Mortems**: Conduct blameless post-mortem reviews after every production outage to identify root causes and assign remediation tasks.

---

## 9. Cost Optimization

Over-provisioning hardware resources wastes money. Apply the following checks:
- **Right-Size VMs**: Monitor CPU and memory usage graphs over a 30-day window. If utilization remains consistently $<10\%$, downgrade the VM resource tier.
- **Lifecycle Rules**: Configure lifecycle policies on storage buckets to automatically delete logs, temporary files, or historical backups after their expiration dates.
- **Managed Service vs. Self-Hosted VPS**:
  - Use managed databases (e.g., AWS RDS) to save database SRE management hours.
  - Use simple VPS instances (e.g., DigitalOcean, Hetzner) for staging environments and small applications to avoid cloud resource bloat.

---

## 10. Legacy Infrastructure

When managing legacy infrastructure that lacks automation or documentation:
1. **Never Make Untracked Changes**: Do not log in to legacy production servers and run modifying configuration commands manually.
2. **Document First**: Record existing server network configurations, environment variables, cron jobs, and database dependencies in Markdown files.
3. **Automate Gradually**: Move manual setups into version-controlled configuration scripts. Automate the build process first, then database backups, and finally the server deployment pipeline.

---

## 11. Decision Matrices

Use these matrices to identify the correct infrastructure decision based on project context.

### Matrix 1: VPS vs. Managed Cloud (AWS / GCP)
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Simple applications, startup MVPs, cost-sensitive staging tracks | **Self-Hosted VPS** | Highly customizable; low upfront costs; simple configuration. |
| Fast growth systems, high-scale security requirements, complex compliance | **Managed Cloud** | Dynamic resource auto-scaling; managed backups and security zones. |

### Matrix 2: Docker vs. Native Deployment
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Multi-stack systems, microservices, multiple hosting platforms | **Docker** | Guarantees runtime consistency; packages system dependencies. |
| Monolithic applications on simple VPS architectures | **Native Deployment**| Avoids registry, container building, and volumes config overhead. |

### Matrix 3: Kubernetes vs. Simple Container Orchestration
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Simple web apps, monolithic platforms, small teams | **Docker Compose** | Simple configuration; runs on a single VPS instance; low maintenance. |
| Large distributed systems, high-scale autoscaling APIs | **Kubernetes / ECS** | Autoscale containers, service discovery, and routing across VM groups. |

### Matrix 4: Redis vs. No Redis
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Shared session storage, rate limit counters, temporary task queues | **Redis** | In-memory speed; prevents database connection bottlenecks. |
| Simple single-node VPS, low-concurrency CRUD apps | **No Redis** | Saves server memory resources and system complexity. |

### Matrix 5: Self-Hosted Database vs. Managed Database
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Staging environments, local testing, cost-sensitive MVPs | **Self-Hosted VPS** | Saves cloud subscription costs. |
| Production business-critical systems, SaaS platforms, financial data | **Managed Database** | Automates patches, backups, and point-in-time recovery configurations. |

### Matrix 6: Vertical Scaling vs. Horizontal Scaling
| Context | Choice | Rationale |
| :--- | :--- | :--- |
| Database systems, stateful applications | **Vertical Scaling**| Simple to scale; maintains consistency without clustering. |
| Stateless web/API application nodes | **Horizontal Scaling**| Redundancy; handles load spikes by adding VM nodes behind a load balancer. |

---

## 12. AI Infrastructure Rules

AI agents suggesting or configuring infrastructure must follow these rules:

1. **Simplicity First**: Never suggest Kubernetes or multi-cluster microservice setups unless explicitly requested by the user.
2. **Never Recommend Raw Sudo Commits**: Do not suggest running raw `sudo` commands directly in bash scripts without documenting their security implications.
3. **No Plaintext Secrets**: Ensure all credential configurations use environment variable references.
4. **Use Pinned Base Configurations**: Pin specific versions for base images and dependency libraries (e.g., `php:8.3-fpm-alpine`, Node version LTS).
5. **No Blind Port Exposures**: Do not expose database, cache, or internal system ports to public network interfaces. Only expose public web ports `80` and `443`.

---

## 13. Infrastructure Review Checklist

Use this checklist during code review to evaluate infrastructure and DevOps changes.

### Security & Access Control
- [ ] Is SSH password authentication disabled (`PasswordAuthentication no`)?
- [ ] Are all database, cache, and internal system ports blocked from public access?
- [ ] Does the Nginx configuration enforce TLS 1.2 or 1.3 (no weak ciphers)?

### Networking & Routing
- [ ] Are DNS records configured correctly?
- [ ] Does Nginx enforce rate limiting zones on API locations?

### Servers & Runtimes
- [ ] Are application processes running under low-privilege system users (no `root` execution)?
- [ ] Do queue workers and background jobs run under process managers (Supervisor, PM2)?

### Backups & Reliability
- [ ] Are database backups scheduled daily and copied to isolated storage?
- [ ] Has restoration testing been performed on backup files?
- [ ] Are monitoring metrics (CPU, RAM, disk space) active?

### Cost Optimization
- [ ] Are VM resource sizes aligned with actual utilization metrics?
- [ ] Are lifecycle rules configured on storage buckets to delete expired logs?

---

## References
- Secure Database Schemas: [06-database-engineering-standard.md](06-database-engineering-standard.md)
- Secure API Gateway Rules: [07-api-engineering-standard.md](07-api-engineering-standard.md)
- Automated CI Pipelines: [13-cicd-and-deployment-standard.md](13-cicd-and-deployment-standard.md)
