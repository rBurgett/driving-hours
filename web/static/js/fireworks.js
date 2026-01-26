// Fireworks celebration animation
(function() {
    const canvas = document.createElement('canvas');
    canvas.style.cssText = 'position:fixed;top:0;left:0;width:100%;height:100%;pointer-events:none;z-index:9999';
    document.body.appendChild(canvas);

    const ctx = canvas.getContext('2d');
    let width, height;
    let particles = [];
    let animationId;

    function resize() {
        width = canvas.width = window.innerWidth;
        height = canvas.height = window.innerHeight;
    }

    resize();
    window.addEventListener('resize', resize);

    class Particle {
        constructor(x, y, color) {
            this.x = x;
            this.y = y;
            this.color = color;
            this.velocity = {
                x: (Math.random() - 0.5) * 8,
                y: (Math.random() - 0.5) * 8 - 2
            };
            this.alpha = 1;
            this.decay = Math.random() * 0.015 + 0.01;
            this.size = Math.random() * 3 + 1;
        }

        update() {
            this.velocity.y += 0.05; // gravity
            this.x += this.velocity.x;
            this.y += this.velocity.y;
            this.alpha -= this.decay;
        }

        draw() {
            ctx.save();
            ctx.globalAlpha = this.alpha;
            ctx.fillStyle = this.color;
            ctx.beginPath();
            ctx.arc(this.x, this.y, this.size, 0, Math.PI * 2);
            ctx.fill();
            ctx.restore();
        }
    }

    function createFirework(x, y) {
        const colors = ['#3b82f6', '#10b981', '#f59e0b', '#ef4444', '#8b5cf6', '#ec4899'];
        const color = colors[Math.floor(Math.random() * colors.length)];
        const particleCount = 50 + Math.floor(Math.random() * 30);

        for (let i = 0; i < particleCount; i++) {
            particles.push(new Particle(x, y, color));
        }
    }

    function animate() {
        ctx.clearRect(0, 0, width, height);

        particles = particles.filter(p => p.alpha > 0);

        particles.forEach(p => {
            p.update();
            p.draw();
        });

        if (particles.length > 0) {
            animationId = requestAnimationFrame(animate);
        } else {
            // Clean up when done
            cancelAnimationFrame(animationId);
            canvas.remove();
        }
    }

    // Launch multiple fireworks
    function launchShow() {
        const launches = 5;
        let launched = 0;

        function launch() {
            if (launched >= launches) return;

            const x = Math.random() * width * 0.6 + width * 0.2;
            const y = Math.random() * height * 0.4 + height * 0.1;
            createFirework(x, y);
            launched++;

            setTimeout(launch, 300 + Math.random() * 400);
        }

        launch();
        animate();
    }

    // Start the show
    launchShow();
})();
