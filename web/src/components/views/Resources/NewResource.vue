<template>
  <modal
    name="newresource"
    class="modal-new-member"
    :draggable="true"
    height=auto
    @before-open="beforeOpen"
    @before-close="beforeClose"
  >
      <div class="box box-group">
        <div class="box-header with-border dark">
          <h3 class="box-title">Add new member</h3>
        </div>
        <div class="box-body">
          <!--Input section-->
          <div class="input-group col-xs-12">
            <span class="input-group-addon">
              <i class="fa fa-tag"></i>
            </span>
            <input
              type="text"
              id="badgeID"
              v-model="member.badgeID"
              v-validate="'required|min:3'"
              class="form-control"
              placeholder="Badge ID"
              name="badgeID"
              :class="{ 'is-invalid': submitted && errors.has('badgeID') }"
           
            />
            <transition name="alert-in" enter-active-class="animated flipInX" leave-active-class="animated flipOutX">
              <div v-if="submitted && errors.has('badgeID')" class="invalid-feedback">{{ errors.first('badgeID') }}</div> 
            </transition>
          </div>
          <div class="input-group col-xs-12">
            <span class="input-group-addon">
              <i class="fa fa-user-plus"></i>
            </span>
            <input
              type="text"
              id="name"
              name="name"
              v-model="member.name"
              class="form-control"
              v-validate="'required|min:3'"
              :class="{ 'is-invalid': submitted && errors.has('name') }"
              placeholder="Username"
            />
            <transition name="alert-in" enter-active-class="animated flipInX" leave-active-class="animated flipOutX">
              <div v-if="submitted && errors.has('name')" class="invalid-feedback">{{ errors.first('name') }}</div>
            </transition>
          </div>
          <div class="input-group col-xs-12">
            <span class="input-group-addon">
              <i class="fa fa-envelope"></i>
            </span>
            <input
              v-validate="'required|email'"
              name="email"
              id="email"
              type="email"
              class="form-control"
              placeholder="Email"
              v-model="member.email"
              :class="{ 'is-invalid': submitted && errors.has('email') }"
            />
            <transition name="alert-in" enter-active-class="animated flipInX" leave-active-class="animated flipOutX">
             <div v-if="submitted && errors.has('email')" class="invalid-feedback">{{ errors.first('email') }}</div>
            </transition>
          </div>
        </div>
        <div class="box-footer">
          <button v-if="isEditting" class="btn-create button-modal pull-right" @click="handleSubmit">Edit</button>
          <button v-else class="btn-create button-modal pull-right" @click="handleSubmit">Invite</button>
          <button class="btn-close button-modal" @click="cancelCreate">Cancel</button>
        </div>
      </div>
  </modal>
</template>
<script>
export default {
  name: 'newresource',
  data() {
    return {
      submitted: false,
      isEditting: false,
      member: {
        id: '',
        badgeID: '',
        name: '',
        email: '',
        project: []
      }
    }
  },
  methods: {
    cancelCreate() {
      this.$modal.hide('newresource')
    },
    handleSubmit() {
      this.submitted = true
      this.$validator.validate().then(valid => {
                if (valid) {
                  if (this.isEditting === false) {
                    this.$store.dispatch('addResource', this.member)
                    this.$modal.hide('newresource')
                  } else {
                    this.$store.dispatch('editResource', this.member)
                    this.$modal.hide('newresource')
                  }
                }
            })
    },
    beforeOpen(event) {
      // check if the dialog open is for edit resources or to add new resources
      if (event.params != null) {
        this.member = event.params.resource
        this.isEditting = true
      } else {
        this.member = {}
        this.member.project = []
      }
    },
    beforeClose() {
      this.isEditting = false
    }
  }
}
</script>
<style scoped>
.input-group {
  margin-top: 15px;
}
.box-title {
  padding: 5px;
  letter-spacing: 1px;
  font-family: "Open Sans", sans-serif;
  font-weight: 400;
  color: #313233;
  text-transform: uppercase;
  transition: 0.1s all;
  font-size: 16px;
}
</style>