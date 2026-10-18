export const drawRotated = (video: CanvasImageSource, canvas: HTMLCanvasElement, context: CanvasRenderingContext2D, degrees: number) => {
  const width = canvas.width;
  const height = canvas.height;
  if (degrees === 90 || degrees === 270) {
    canvas.width = height;
    canvas.height = width;
  }
  if (degrees === 90) {
    context.translate(height, 0);
  }
  else if (degrees === 180) {
    context.translate(width, height);
  }
  else if (degrees === 270) {
    context.translate(0, width);
  }
  else {
    context.translate(0, 0);
  }
  context.rotate(degrees * Math.PI / 180);

  // the frame whole, from the corner the rotation moved it to
  context.drawImage(video, 0, 0, width, height);
}

/** screenRotation is how much the screen is rotated (0, 90, 180, 270): a
 * picture taken is turned back by it. 0 where the browser does not say. */
export const screenRotation = (): number => {
  const angle = (typeof screen !== 'undefined' && screen.orientation?.angle) ??
    ((typeof window !== 'undefined' && (window as any).orientation) || 0)
  const a = ((Number(angle) % 360) + 360) % 360
  return [0, 90, 180, 270].includes(a) ? a : 0
}
