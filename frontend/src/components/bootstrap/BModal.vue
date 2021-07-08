<template>
  <div>
    <button
      type="button"
      class="button btn-primary"
      data-bs-toggle="modal"
      :data-bs-target="`#${targetId}`"
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
            <template v-if="modalType === 'ok-cancel'">
              <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Close</button>
              <button type="button" class="btn btn-success" data-bs-dismiss="modal" @click="proceeded">Save</button>
            </template>
            <template v-if="modalType === 'confirm'">
              <button type="button" class="btn btn-success" data-bs-dismiss="modal" @click="proceeded">Save</button>
            </template>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, SetupContext } from 'vue'

import BButton from '@/components/bootstrap/BButton.vue'

type ModalType = 'close' | 'ok-cancel' | 'confirm'

type Props = {
  buttonTitle: string
  modalTitle: string
  modalType: ModalType
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
    modalTitle: {
      type: String,
      default: '',
    },
    modalType: {
      type: String as () => ModalType,
      default: 'close',
    }
  },
  setup(_: Props, context: SetupContext) {
    const targetId = 'tg' + Math.random().toString(32).substring(2)

    const proceeded = (() => {
      context.emit('proceeded')
    })

    return {
      targetId,
      proceeded,
    }
  },
})
</script>

<style scoped>
.modal-title {
  padding-top: 2px;
}
</style>