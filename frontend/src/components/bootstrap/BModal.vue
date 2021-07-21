<template>
  <div>
    <button
      type="button"
      :class="`btn ${openButtonVariant}`"
      data-bs-toggle="modal"
      :data-bs-target="`#${targetId}`"
      @click="openModal"
    >
      {{ buttonTitle }}
    </button>

    <div class="modal fade" :id="targetId">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header" v-if="modalTitle">
            <h4 class="modal-title" :id="`${targetId}-title`">{{ modalTitle }}</h4>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <slot />
          </div>
          <div class="modal-footer">
            <template v-if="modalType === 'close'">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
            </template>
            <template v-if="modalType === 'confirm-cancel'">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">{{ cancelTitle }}</button>
              <button type="button" :class="`btn ${confirmButtonVariant}`" data-bs-dismiss="modal" @click="proceeded">{{ confirmTitle }}</button>
            </template>
            <template v-if="modalType === 'confirm'">
              <button type="button" :class="`btn ${confirmButtonVariant}`" data-bs-dismiss="modal" @click="proceeded">{{ confirmTitle }}</button>
            </template>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { computed, defineComponent, SetupContext } from 'vue'

import BButton from '@/components/bootstrap/BButton.vue'

type ModalType = 'close' | 'confirm-cancel' | 'confirm'

type Props = {
  buttonTitle: string
  buttonVariant: string
  modalTitle: string
  modalType: ModalType
  confirmTitle: string
  confirmVariant: string
  cancelTitle: string
}

export default defineComponent({
  name: 'BModal',
  components: {
    BButton,
  },
  props: {
    buttonTitle: {
      type: String,
      default: 'Open modal',
    },
    buttonVariant: {
      type: String,
      default: 'primary',
    },
    modalTitle: {
      type: String,
      default: '',
    },
    modalType: {
      type: String as () => ModalType,
      default: 'close',
    },
    confirmTitle: {
      type: String,
      default: 'Save',
    },
    confirmVariant: {
      type: String,
      default: 'success',
    },
    cancelTitle: {
      type: String,
      default: 'Close',
    },
  },
  emits: ['visible', 'proceeded'],
  setup(props: Props, context: SetupContext) {
    const targetId = 'tg' + Math.random().toString(32).substring(2)

    const openButtonVariant = computed(() => {
      if (props.buttonVariant) {
        return `btn-${props.buttonVariant}`
      }
      return 'btn-primary'
    })

    const confirmButtonVariant = computed(() => {
      if (props.confirmVariant) {
        return `btn-${props.confirmVariant}`
      }
      return 'btn-primary'
    })

    const proceeded = (() => {
      context.emit('proceeded')
    })

    const openModal = (() => {
      context.emit('visible')
    })

    return {
      targetId,
      proceeded,
      openButtonVariant,
      confirmButtonVariant,
      openModal,
    }
  },
})
</script>

<style scoped>
.modal-title {
  padding-top: 2px;
}
</style>