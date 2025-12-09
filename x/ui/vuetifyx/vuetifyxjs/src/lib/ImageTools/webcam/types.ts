export class Picture {
  data: Blob | null = null
  url: String = ''

  clear() {
    this.data = null
    this.url = ''
  }

  set(data: Blob, url: String) {
    this.data = data
    this.url = url
  }
}
