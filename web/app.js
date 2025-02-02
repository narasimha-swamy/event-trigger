let currentTriggerId = null;
let currentTriggerType = null;

document.addEventListener('DOMContentLoaded', () => {
    loadTriggers();
    loadEventLogs();
});

document.getElementById('showArchived').addEventListener('change', () => {
    loadEventLogs();
});

// Add this toast function at the top of app.js
function showToast(message, type = 'success') {
    const toastContainer = document.createElement('div');
    toastContainer.innerHTML = `
        <div class="toast align-items-center text-white bg-${type === 'success' ? 'success' : 'danger'} border-0" role="alert">
            <div class="d-flex">
                <div class="toast-body">${message}</div>
                <button type="button" class="btn-close btn-close-white me-2 m-auto" data-bs-dismiss="toast"></button>
            </div>
        </div>
    `;
    
    document.body.appendChild(toastContainer);
    const toast = new bootstrap.Toast(toastContainer.querySelector('.toast'));
    toast.show();
    
    // Auto-remove after animation
    setTimeout(() => toastContainer.remove(), 3000);
}
async function loadTriggers() {
    try {
        const response = await fetch('/api/triggers');
        const triggers = await response.json();
        renderTriggers(triggers);
    } catch (error) {
        showError('Failed to load triggers');
    }
}

function loadEventLogs() {
    const showArchived = document.getElementById('showArchived').checked;
    const url = `/api/events${showArchived ? '?archived=true' : ''}`;
    
    fetch(url)
        .then(response => response.json())
        .then(events => renderEventLogs(events))
        .catch(error => showError('Failed to load event logs'));
}


// Update the renderTriggers function to ensure proper event binding
function renderTriggers(triggers) {
    const container = document.getElementById('triggerList');
    container.innerHTML = triggers.map(trigger => `
        <div class="card mb-2">
            <div class="card-body">
                <h5>${trigger.type} Trigger</h5>
                <p class="text-muted small">ID: ${trigger.id}</p>
                ${trigger.type === 'scheduled' ? `
                    <p>Cron: <code>${trigger.cron_expression}</code></p>
                    <p>Next Run: ${new Date(trigger.next_run).toLocaleString()}</p>
                ` : ''}
                <div class="mt-2">
                    <button class="btn btn-sm btn-danger" onclick="deleteTrigger('${trigger.id}')">
                        Delete
                    </button>
                </div>
            </div>
        </div>
    `).join('');
    
    // Reinitialize modal handlers
    initializeTestModal();
}

  function initializeTestModal() {
    const modal = new bootstrap.Modal('#testTriggerModal');
    const payloadSection = document.getElementById('payloadSection');
    const scheduledNotice = document.getElementById('scheduledNotice');
    
    document.querySelectorAll('.test-trigger').forEach(button => {
      button.addEventListener('click', (e) => {
        currentTriggerId = e.target.dataset.triggerId;
        currentTriggerType = e.target.dataset.triggerType;
        
        if (currentTriggerType === 'api') {
          document.getElementById('modalTitle').textContent = 'Test API Trigger';
          payloadSection.style.display = 'block';
          scheduledNotice.style.display = 'none';
        } else {
          document.getElementById('modalTitle').textContent = 'Test Scheduled Trigger';
          payloadSection.style.display = 'none';
          scheduledNotice.style.display = 'block';
        }
        
        modal.show();
      });
    });
  
    document.getElementById('confirmTest').addEventListener('click', async () => {
        try {
            const modal = bootstrap.Modal.getInstance('#testTriggerModal');
            let payload = null;
            
            if (currentTriggerType === 'api') {
                const payloadInput = document.getElementById('testPayload').value;
                try {
                    payload = JSON.parse(payloadInput);
                } catch (e) {
                    throw new Error('Invalid JSON format');
                }
            }
    
            const response = await fetch(`/api/triggers/${currentTriggerId}/test`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: payload ? JSON.stringify(payload) : null
            });
    
            const result = await response.json();
            
            if (!response.ok) {
                throw new Error(result.error || 'Test failed');
            }
    
            // Show success toast
            showToast('Trigger tested successfully!', 'success');
            modal.hide();
            loadEventLogs();
        } catch (error) {
            // Show error in modal
            document.getElementById('payloadError').textContent = error.message;
            document.getElementById('payloadSection').classList.add('was-validated');
            console.error('Test failed:', error);
        }
    });
  }

// Render event logs
function renderEventLogs(events) {
    const container = document.getElementById('eventLogs');
    container.innerHTML = events.map(event => `
        <div class="log-entry ${event.status === 'archived' ? 'text-muted' : ''}">
            <div>${new Date(event.triggered_at).toLocaleString()}</div>
            <div>Trigger: ${event.trigger_id}</div>
            ${event.payload ? `<pre>${JSON.stringify(event.payload, null, 2)}</pre>` : ''}
            ${event.status === 'archived' ? '<div class="badge bg-secondary">Archived</div>' : ''}
        </div>
    `).join('');
}

async function deleteTrigger(triggerId) {
    try {
        await fetch(`/api/triggers/${triggerId}`, { method: 'DELETE' });
        loadTriggers();
    } catch (error) {
        showError('Failed to delete trigger');
    }
}

function showCreateForm(type) {
    const form = type === 'scheduled' ? `
        <form onsubmit="createScheduledTrigger(event)">
            <div class="mb-3">
                <label>Cron Expression</label>
                <input type="text" class="form-control" required 
                    placeholder="* * * * *">
            </div>
            <div class="form-check mb-3">
                <input type="checkbox" class="form-check-input" id="recurring">
                <label class="form-check-label" for="recurring">Recurring</label>
            </div>
            <button type="submit" class="btn btn-success">Create</button>
        </form>
    ` : `
        <form onsubmit="createAPITrigger(event)">
            <div class="mb-3">
                <label>Payload Schema (JSON)</label>
                <textarea class="form-control" rows="3" required
                    placeholder='{"key":"value"}'></textarea>
            </div>
            <button type="submit" class="btn btn-success">Create</button>
        </form>
    `;
    
    document.getElementById('createForm').innerHTML = form;
}

async function createScheduledTrigger(event) {
    event.preventDefault();
    // const form = event.target;
    // const feedbackDiv = form.parentElement.querySelector('#scheduledFeedback');
    
    // Get form data
    const cronExpression = event.target[0].value;
    const isRecurring = event.target[1].checked;
    
    try {
        // Create trigger
        const response = await fetch('/api/triggers', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                type: 'scheduled',
                cron_expression: cronExpression,
                is_recurring: isRecurring,
                is_active: true
            })
        });

        if (!response.ok) {
            throw new Error('Failed to create trigger');
        }

        // Refresh trigger list
        loadTriggers();
        
        // Clear form
        event.target.reset();
        
        // Show success message
        showSuccess('Scheduled trigger created successfully!');
    } catch (error) {
        showError(`Failed to create trigger: ${error.message}`);
    }
    
}

function showError(message) {
    const alert = document.createElement('div');
    alert.className = 'alert alert-danger';
    alert.textContent = message;
    document.querySelector('.container').prepend(alert);
    setTimeout(() => alert.remove(), 3000);
}

function showFeedback(message, type, containerId) {
    const feedbackDiv = document.getElementById(containerId);
    feedbackDiv.innerHTML = `
        <div class="alert alert-${type} alert-dismissible fade show" role="alert">
            ${message}
            <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
        </div>
    `;
}

async function createAPITrigger(event) {
    event.preventDefault();
    const form = event.target;
    const payloadInput = form.querySelector('textarea');
    const feedbackDiv = form.parentElement.querySelector('#apiFeedback');
    
    try {
        const payload = JSON.parse(payloadInput.value);
        const response = await fetch('/api/triggers', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
                type: 'api',
                api_payload: payload,
                is_active: true
            })
        });

        if (!response.ok) throw new Error(await response.text());
        
        payloadInput.value = '';
        showFeedback('API trigger created successfully!', 'success', 'apiFeedback');
        loadTriggers();
    } catch (error) {
        showFeedback(error.message, 'danger', 'apiFeedback');
    }
}

function showSuccess(message) {
    const alert = document.createElement('div');
    alert.className = 'alert alert-success';
    alert.textContent = message;
    document.querySelector('.container').prepend(alert);
    setTimeout(() => alert.remove(), 3000);
}

async function testTrigger(triggerId) {
    try {
        let payload = null;
        if (confirm("Is this an API trigger? Provide payload if needed.")) {
            const payloadSchema = prompt("Enter payload (JSON):");
            payload = JSON.parse(payloadSchema);
        }

        const response = await fetch(`/api/triggers/${triggerId}/test`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: payload ? JSON.stringify(payload) : null
        });

        if (!response.ok) {
            throw new Error('Failed to test trigger');
        }

        showSuccess('Trigger tested successfully!');
        loadEventLogs(); // Refresh event logs
    } catch (error) {
        showError(`Failed to test trigger: ${error.message}`);
    }
}

// Add to app.js
async function testTriggerWithPayload(event) {
    event.preventDefault();
    const feedbackDiv = document.getElementById('testFeedback');
    const triggerId = document.getElementById('testTriggerId').value;
    const payloadInput = document.getElementById('testPayload').value;
    feedbackDiv.innerHTML = '';

    try {
        let payload = {};
        if (payloadInput) {
            payload = JSON.parse(payloadInput);
        }

        const response = await fetch(`/api/triggers/${triggerId}/test`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(payload)
        });

        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error || 'Test failed');
        }

        showTestFeedback('Trigger tested successfully! Event logged.', 'success');
        loadEventLogs(); // Refresh event logs
        
    } catch (error) {
        showTestFeedback(`Test failed: ${error.message}`, 'danger');
    }
}

function showTestFeedback(message, type) {
    const feedbackDiv = document.getElementById('testFeedback');
    feedbackDiv.innerHTML = `
        <div class="alert alert-${type} alert-dismissible fade show mt-2" role="alert">
            ${message}
            <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close"></button>
        </div>
    `;
}

// Update modal handling
function initializeTestModal() {
    const modal = new bootstrap.Modal('#testTriggerModal');
    const scheduledUI = document.getElementById('scheduledTestUI');
    const apiUI = document.getElementById('apiTestUI');

    document.querySelectorAll('.test-trigger').forEach(button => {
        button.addEventListener('click', (e) => {
            currentTriggerId = e.target.dataset.triggerId;
            currentTriggerType = e.target.dataset.triggerType;
            
            scheduledUI.style.display = 'none';
            apiUI.style.display = 'none';

            if (currentTriggerType === 'api') {
                apiUI.style.display = 'block';
            } else {
                scheduledUI.style.display = 'block';
            }
            
            modal.show();
        });
    });

    document.getElementById('confirmTest').addEventListener('click', async () => {
        try {
            let body = null;
            
            if (currentTriggerType === 'scheduled') {
                const delay = document.getElementById('testDelayMinutes').value;
                body = { delay_minutes: parseInt(delay) };
            } else {
                // Existing API payload handling
            }

            const response = await fetch(`/api/triggers/${currentTriggerId}/test`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify(body)
            });

            if (!response.ok) throw new Error(await response.text());
            
            showToast(`Test scheduled successfully!`, 'success');
            modal.hide();
        } catch (error) {
            showToast(error.message, 'danger');
        }
    });
}

function showTestModal(triggerId, triggerType) {
    currentTriggerId = triggerId;
    currentTriggerType = triggerType;
    const modal = new bootstrap.Modal('#testTriggerModal');
    modal.show();
}

async function handleTestTrigger() {
    const paramsInput = document.getElementById('testParams').value;
    const errorDiv = document.getElementById('testParamsError');
    
    try {
        const params = JSON.parse(paramsInput);
        
        const response = await fetch(`/api/triggers/${currentTriggerId}/test`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify(params)
        });

        const result = await response.json();
        
        if (!response.ok) throw new Error(result.error);
        
        showToast(`Test scheduled: ${result.message}`, 'success');
        bootstrap.Modal.getInstance('#testTriggerModal').hide();
        loadEventLogs();
    } catch (error) {
        errorDiv.textContent = error.message;
        document.getElementById('testParams').classList.add('is-invalid');
    }
}