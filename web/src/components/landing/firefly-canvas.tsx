import { useEffect, useRef } from 'react';
import * as THREE from 'three';

/*
 * Fireflies drifting through the night garden. Fixed behind the whole
 * landing page, additively blended so every spore glows against the dark.
 * Lazy-loaded so three.js never blocks the first paint.
 */

const PALETTE = ['#d3f36b', '#9fdc6c', '#5f9e6e', '#e7c66b', '#f2eedd'];

function makeSprite(): THREE.CanvasTexture {
	const canvas = document.createElement('canvas');
	canvas.width = canvas.height = 64;
	const ctx = canvas.getContext('2d')!;
	const gradient = ctx.createRadialGradient(32, 32, 0, 32, 32, 32);
	gradient.addColorStop(0, 'rgba(255,255,255,1)');
	gradient.addColorStop(0.3, 'rgba(255,255,255,0.6)');
	gradient.addColorStop(1, 'rgba(255,255,255,0)');
	ctx.fillStyle = gradient;
	ctx.fillRect(0, 0, 64, 64);
	return new THREE.CanvasTexture(canvas);
}

interface Layer {
	points: THREE.Points;
	geometry: THREE.BufferGeometry;
	material: THREE.PointsMaterial;
	baseOpacity: number;
	twinklePhase: number;
	speeds: Float32Array;
	phases: Float32Array;
	sways: Float32Array;
}

function makeLayer(
	count: number,
	size: number,
	opacity: number,
	sprite: THREE.Texture
): Layer {
	const positions = new Float32Array(count * 3);
	const colors = new Float32Array(count * 3);
	const speeds = new Float32Array(count);
	const phases = new Float32Array(count);
	const sways = new Float32Array(count);
	const color = new THREE.Color();

	for (let i = 0; i < count; i++) {
		positions[i * 3] = (Math.random() - 0.5) * 56;
		positions[i * 3 + 1] = (Math.random() - 0.5) * 34;
		positions[i * 3 + 2] = (Math.random() - 0.5) * 14;
		color.set(PALETTE[Math.floor(Math.random() * PALETTE.length)]);
		colors[i * 3] = color.r;
		colors[i * 3 + 1] = color.g;
		colors[i * 3 + 2] = color.b;
		speeds[i] = 0.2 + Math.random() * 0.5;
		phases[i] = Math.random() * Math.PI * 2;
		sways[i] = 0.15 + Math.random() * 0.55;
	}

	const geometry = new THREE.BufferGeometry();
	geometry.setAttribute('position', new THREE.BufferAttribute(positions, 3));
	geometry.setAttribute('color', new THREE.BufferAttribute(colors, 3));

	const material = new THREE.PointsMaterial({
		size,
		map: sprite,
		vertexColors: true,
		transparent: true,
		opacity,
		depthWrite: false,
		sizeAttenuation: true,
		blending: THREE.AdditiveBlending,
	});

	return {
		points: new THREE.Points(geometry, material),
		geometry,
		material,
		baseOpacity: opacity,
		twinklePhase: Math.random() * Math.PI * 2,
		speeds,
		phases,
		sways,
	};
}

export default function FireflyCanvas() {
	const mountRef = useRef<HTMLDivElement>(null);

	useEffect(() => {
		const mount = mountRef.current;
		if (!mount) return;

		const reducedMotion = window.matchMedia(
			'(prefers-reduced-motion: reduce)'
		).matches;
		const coarsePointer = window.matchMedia('(pointer: coarse)').matches;

		const scene = new THREE.Scene();
		const camera = new THREE.PerspectiveCamera(60, 1, 0.1, 100);
		camera.position.z = 24;

		const renderer = new THREE.WebGLRenderer({
			alpha: true,
			antialias: true,
			powerPreference: 'low-power',
		});
		renderer.setClearColor(0x000000, 0);
		renderer.domElement.style.opacity = '0';
		renderer.domElement.style.transition = 'opacity 1.8s ease';
		mount.appendChild(renderer.domElement);

		const sprite = makeSprite();
		const layers: Layer[] = [
			makeLayer(coarsePointer ? 90 : 190, 0.3, 0.4, sprite), // fine spores
			makeLayer(coarsePointer ? 40 : 85, 0.85, 0.6, sprite), // fireflies
		];
		layers.forEach((layer) => scene.add(layer.points));

		const pointer = { x: 0, y: 0 };
		const onPointerMove = (event: PointerEvent) => {
			pointer.x = (event.clientX / window.innerWidth) * 2 - 1;
			pointer.y = (event.clientY / window.innerHeight) * 2 - 1;
		};
		if (!coarsePointer) {
			window.addEventListener('pointermove', onPointerMove, {
				passive: true,
			});
		}

		const resize = () => {
			const { clientWidth, clientHeight } = mount;
			camera.aspect = clientWidth / Math.max(clientHeight, 1);
			camera.updateProjectionMatrix();
			renderer.setPixelRatio(Math.min(window.devicePixelRatio, 2));
			renderer.setSize(clientWidth, clientHeight);
		};
		resize();
		const resizeObserver = new ResizeObserver(resize);
		resizeObserver.observe(mount);

		let raf = 0;
		const clock = new THREE.Clock();

		const render = () => {
			// getElapsedTime() would internally consume getDelta(), so read
			// the delta first and the accumulated time from the property.
			const delta = Math.min(clock.getDelta(), 0.05);
			const elapsed = clock.elapsedTime;

			for (const layer of layers) {
				const positions = layer.geometry.getAttribute(
					'position'
				) as THREE.BufferAttribute;
				const array = positions.array as Float32Array;
				for (let i = 0; i < layer.speeds.length; i++) {
					array[i * 3 + 1] += layer.speeds[i] * delta;
					array[i * 3] +=
						Math.sin(elapsed * 0.5 + layer.phases[i]) *
						layer.sways[i] *
						delta;
					if (array[i * 3 + 1] > 18) array[i * 3 + 1] = -18;
				}
				positions.needsUpdate = true;
				layer.material.opacity =
					layer.baseOpacity *
					(0.82 + Math.sin(elapsed * 1.3 + layer.twinklePhase) * 0.18);
			}

			// Pointer parallax plus a whisper of scroll drift.
			const scrollDrift =
				Math.sin(window.scrollY * 0.0004) * 1.2;
			camera.position.x +=
				(pointer.x * 2.2 + scrollDrift - camera.position.x) * 0.04;
			camera.position.y += (-pointer.y * 1.4 - camera.position.y) * 0.04;
			camera.lookAt(scene.position);
			renderer.render(scene, camera);
		};

		const loop = () => {
			render();
			raf = requestAnimationFrame(loop);
		};

		const onVisibility = () => {
			cancelAnimationFrame(raf);
			if (!document.hidden && !reducedMotion) {
				clock.getDelta();
				raf = requestAnimationFrame(loop);
			}
		};

		if (reducedMotion) {
			render(); // single static frame
		} else {
			raf = requestAnimationFrame(loop);
			document.addEventListener('visibilitychange', onVisibility);
		}

		requestAnimationFrame(() => {
			renderer.domElement.style.opacity = '1';
		});

		return () => {
			cancelAnimationFrame(raf);
			document.removeEventListener('visibilitychange', onVisibility);
			window.removeEventListener('pointermove', onPointerMove);
			resizeObserver.disconnect();
			layers.forEach((layer) => {
				layer.geometry.dispose();
				layer.material.dispose();
			});
			sprite.dispose();
			renderer.dispose();
			renderer.domElement.remove();
		};
	}, []);

	return (
		<div
			ref={mountRef}
			aria-hidden
			className="pointer-events-none fixed inset-0 z-0"
		/>
	);
}
