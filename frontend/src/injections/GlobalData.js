/**
 * @type {import('vue').InjectionKey<{
 *   settings: import('vue').Ref<Record<string, any>>,
 *   refreshSettings: () => void,
 *   setLoading: (loading: boolean) => void,
 *   message: {
 *     success: (content: string) => void,
 *     error: (content: string) => void,
 *   },
 * }>}
 */
export default 'GlobalData'