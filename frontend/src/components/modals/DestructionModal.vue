<template>
  <BModal
    :button-title="buttonTitle"
    button-variant="danger"
    modal-type="confirm-cancel"
    :modal-title="modalTitle"
    confirm-variant="danger"
    confirm-title="Yes, Delete"
    cancel-title="Cancel"

    @proceeded="proceedDelete"
  >
    <p>Are you sure to delete this?</p>
  </BModal>
</template>

<script lang="ts">
import { ref, defineComponent, SetupContext } from 'vue'
// import axios from 'axios'
import BModal from '@/components/bootstrap/BModal.vue'
// import BInput from '@/components/bootstrap/BInput.vue'


export default defineComponent({
  name: 'DestructionModal',
  components: {
    BModal,
    // BInput,
  },
  props: {
    buttonTitle: {
      type: String,
      default: 'Delete',
    },
    modalTitle: {
      type: String,
      default: 'Confirm Deletion',
    }
  },
  emits: ['confirmDelete'],
  setup(props, context: SetupContext) {
    const name = ref('')
    const description = ref('')

    const proceedDelete = (() => {
      context.emit('confirmDelete')
    })

    return {
      name,
      description,

      proceedDelete,
    }
  },
})
</script>
