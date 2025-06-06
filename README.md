# 📬 SalesForge Email Sequencing System

A robust, scalable, and reliable email sequencing and delivery system with intelligent scheduling, mailbox quota handling, and retry mechanisms.

---

## 🚀 Key Concepts and Components

### 🧱 Data Model Enhancements

| Table             | Purpose                                                                 |
|------------------|-------------------------------------------------------------------------|
| `sequences`       | Defines sending window and linked `sequence_steps`.                    |
| `sequence_steps`  | Defines wait days before sending the next step to a contact.           |
| `contacts`        | Each contact can be associated with one or more sequences.             |
| `mailboxes`       | Stores per-mailbox sending quota; tracks `quota_used` per day.         |
| `scheduled_emails`| Core queue with `timestamp`, `mailbox_id`, `contact_id`, `step_id`, and `status`. |

---

## 🕓 How Email Scheduling Works

### Step-by-Step Flow

#### ✅ Nightly Scheduler Service (Runs at midnight)

- Iterates over all active sequences.
- For each contact:
  - Checks last step and waiting period.
  - Determines next eligible step.
  - Calculates next send time within sending window, spaced evenly.
  - Assigns a mailbox with remaining quota.
  - Inserts entry into `scheduled_emails`.

#### 📊 Quota & Distribution Logic

- **Example**: Sending window = 10 hours = 600 mins.
- Each mailbox = 30 emails/day → **1 email every 20 mins**.
- Scheduler staggers email timings **per mailbox** to avoid burst sends.

---

### 📤 Dispatcher Service

- Continuously monitors `scheduled_emails`:
  - `scheduled_at <= now`
  - `status = 'pending'`
- Pushes jobs to queue (RabbitMQ, Kafka, or Postgres-based PGQ).
- Ensures high throughput and batching for efficiency.

---

### 🛠 Email Worker

- Pulls job from queue.
- Sends email via SMTP from assigned mailbox.
- Marks job `sent`, or `failed` with retry counter.
- Retries with exponential backoff on transient errors.

---

## 🔁 Retry Strategy

- **Max Retries**: 3
- **Backoff Schedule**: 5 min → 15 min → 30 min
- On failure beyond retry threshold: mark as permanently failed.

---

## ⚖️ Horizontal Scaling Strategy

| Component         | Scaling Approach                                                   |
|------------------|---------------------------------------------------------------------|
| Scheduler         | Stateless, one instance with locking (e.g., advisory DB locks).     |
| Dispatcher        | Stateless, can scale with job queue consumers.                     |
| Email Workers     | Horizontally scalable (Kubernetes pods, containers, etc.).         |
| Mailboxes         | Enforced via transactional row updates to `quota_used`.            |

---

## 💾 Data Persistence

- **Primary DB**: PostgreSQL (or CockroachDB for multi-region support).
- **Redis (optional)**: Caching, locking, or mailbox quota reference.

### Schema Highlights

- `mailboxes(id, email, quota_used, reset_time)`
- `scheduled_emails(id, step_id, contact_id, scheduled_at, status)`
- `contacts`, `sequences`, `sequence_steps`

---

## 🔒 Reliability & Idempotency

- Email send operations are **idempotent**.
- Mailbox `quota_used` updates occur inside **atomic DB transactions**.
- Locking per-mailbox enforced in DB to avoid quota races.

---

## 📈 Observability

- Logs per email job (success, failure, retries).
- Metrics dashboard for:
  - Emails sent per mailbox
  - Failure rate
  - Quota exhaustion alerts
- Alerts if mailbox quota consumed early in the day.

---

## 📊 Example Scenario

**Sequence**: “Welcome Campaign”  
**Steps**:  
- Step 1: Immediate  
- Step 2: Wait 2 days  
- Step 3: Wait 1 day  

**Sending Window**: 9 AM – 5 PM (8 hours = 480 mins)  
**Mailboxes**: 3 (30 emails/day each)  
**Total Emails**: 90

→ Scheduler assigns **30 emails per mailbox**, staggered **every 16 minutes** throughout the window.

---

## 🧪 Future Enhancements

- Priority queues (marketing vs. transactional).
- AI-based mailbox rotation based on delivery performance.
- UI dashboard for visualizing quota & schedule.
- Timezone-aware delivery logic (per recipient location).

---

