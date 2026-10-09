<script setup lang="ts">
/** 受控开关。用 button 而不是 span，尺寸和可聚焦性都有保证。 */
const props = withDefaults(
  defineProps<{
    modelValue: boolean
    disabled?: boolean
    label?: string
  }>(),
  { disabled: false, label: '' },
)
const emit = defineEmits<{ 'update:modelValue': [boolean] }>()

function toggle() {
  if (props.disabled) return
  emit('update:modelValue', !props.modelValue)
}
</script>

<template>
  <button
    type="button"
    role="switch"
    class="dh-switch"
    :data-on="modelValue ? 'true' : 'false'"
    :disabled="disabled"
    :aria-checked="modelValue"
    :aria-label="label || undefined"
    @click="toggle"
  />
</template>
