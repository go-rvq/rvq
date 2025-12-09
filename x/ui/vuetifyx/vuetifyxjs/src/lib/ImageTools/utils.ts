export const fileToURL = (f: File): Promise<string> => {
  return new Promise<string>((resolve) => {
    const r = new FileReader()
    r.onloadend = () => {
      resolve(r.result as string)
    }
    r.readAsDataURL(f)
  })
}
