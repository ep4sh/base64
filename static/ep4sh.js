  document.addEventListener('htmx:afterOnLoad', function(evt) {
    const btn = document.getElementById('process-btn');
    if (!btn) return;

    btn.classList.remove('htmx-request', 'htmx-error');

    if (evt.detail.xhr.status >= 400) {
      btn.classList.add('htmx-error');
    }
  });

  document.getElementById('process-btn')?.addEventListener('click', function() {
    this.classList.remove('htmx-error');
  });
