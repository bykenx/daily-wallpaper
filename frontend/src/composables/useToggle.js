import { ref } from 'vue'

/**
 * 
 * @returns {[import('vue').Ref, () => void]}
 */
export default function useToggle() {
  const value = ref(false)
  return [value, () => value.value = !value.value]
}