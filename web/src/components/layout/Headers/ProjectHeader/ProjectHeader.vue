<template>
  <header class="main-header">
          <add-sum-task-modal pivotX="1.0"></add-sum-task-modal>
          <add-task-modal></add-task-modal>
          <task-modal :exceptionDays="exceptionDays"></task-modal>
          <break-task-modal :exceptionDays="exceptionDays"></break-task-modal>
          <assign-member-modal :users="users"></assign-member-modal>
          <save-modal :projectId="id"></save-modal>
          <filter-modal></filter-modal>
          <setting-modal :project="project"></setting-modal>
      <div class="navbar-custom-menu">
        <ul class="nav navbar-nav">
          <status></status>
          <button-text></button-text>
          <button-add-task></button-add-task>
          <button-add-member :id="id" :users="users"></button-add-member>
          <button-history :id="id"></button-history>
          <button-filter></button-filter>
          <button-setting></button-setting>
          <button-save></button-save>
        </ul>
      </div>
      <v-dialog/>
  </header>
</template>

<script>
// Buttons
import ButtonHistory from './Buttons/ButtonHistory'
import ButtonSave from './Buttons/ButtonSave'
import ButtonSetting from './Buttons/ButtonSetting'
import ButtonFilter from './Buttons/ButtonFilter'
import ButtonAddMember from './Buttons/ButtonAddMember'
import ButtonAddTask from './Buttons/ButtonAddTask'
import ButtonText from './Buttons/ButtonText'
// Modals
import SaveModal from './Modals/SaveModal'
import FilterModal from './Modals/FilterModal'
import SettingModal from './Modals/SettingModal'
import TaskModal from './Modals/TaskModal'
import AddSumTaskModal from './Modals/AddSumTaskModal'
import AddTaskModal from './Modals/AddTaskModal'
import AssignMemberModal from './Modals/AssignMemberModal'
import BreakTaskModal from './Modals/BreakTaskModal'
import Status from './Modals/Status'

export default {
  name: 'ProjectHeader',
  components: {
    ButtonHistory,
    ButtonSave,
    ButtonSetting,
    ButtonFilter,
    ButtonAddMember,
    ButtonAddTask,
    AddSumTaskModal,
    SettingModal,
    FilterModal,
    TaskModal,
    AddTaskModal,
    BreakTaskModal,
    AssignMemberModal,
    ButtonText,
    Status,
    SaveModal
  },
  props: ['id', 'users', 'exceptionDays', 'project']
}
</script>

<style scoped>
.header-name {
  width: 50%;
  display: fixed;
}
.navbar-nav {
  font-size: 18px;
}

</style>
