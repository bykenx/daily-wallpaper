import { ref } from 'vue'

export default function useToggle() {
  const value = ref(false)
  return [value, () => value.value = !value.value]
}