import { useEffect, useRef } from "react";

/**
 * Depth layer behind the app. Orbs sit at different translateZ depths inside a
 * shared perspective, so one pointer offset projects into real parallax.
 */
export function AmbientBackground() {
  const layerRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const layer = layerRef.current;
    if (!layer) return;
    if (window.matchMedia("(prefers-reduced-motion: reduce)").matches) return;
    if (window.matchMedia("(pointer: coarse)").matches) return;

    let frame = 0;
    const handlePointerMove = (event: PointerEvent) => {
      cancelAnimationFrame(frame);
      frame = requestAnimationFrame(() => {
        const x = (event.clientX / window.innerWidth - 0.5) * 24;
        const y = (event.clientY / window.innerHeight - 0.5) * 24;
        layer.style.setProperty("--drift-x", `${x}px`);
        layer.style.setProperty("--drift-y", `${y}px`);
      });
    };

    window.addEventListener("pointermove", handlePointerMove, {
      passive: true,
    });
    return () => {
      window.removeEventListener("pointermove", handlePointerMove);
      cancelAnimationFrame(frame);
    };
  }, []);

  return (
    <div
      ref={layerRef}
      aria-hidden="true"
      className="pointer-events-none fixed inset-0 overflow-hidden [perspective:1000px]"
    >
      <div className="absolute -top-40 -left-32 h-[30rem] w-[30rem] rounded-full bg-mint-500/20 blur-[110px] [transform:translate3d(var(--drift-x,0px),var(--drift-y,0px),140px)]" />
      <div className="absolute top-1/4 -right-44 h-[34rem] w-[34rem] rounded-full bg-mint-600/16 blur-[120px] [transform:translate3d(var(--drift-x,0px),var(--drift-y,0px),40px)]" />
      <div className="absolute -bottom-52 left-1/5 h-[28rem] w-[28rem] rounded-full bg-ink-800/70 blur-[100px] [transform:translate3d(var(--drift-x,0px),var(--drift-y,0px),-60px)]" />
      <div className="grain absolute inset-0 opacity-15 mix-blend-soft-light" />
    </div>
  );
}
