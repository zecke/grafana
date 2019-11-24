import { actionCreatorFactory } from '../../../../core/redux';

export interface PanelEditorInitCompleted {
  activeTab: PanelEditorTabIds;
  tabs: PanelEditorTab[];
}

export const panelEditorCleanUp = actionCreatorFactory('PANEL_EDITOR_CLEAN_UP').create();

export const panelEditorInitCompleted = actionCreatorFactory<PanelEditorInitCompleted>(
  'PANEL_EDITOR_INIT_COMPLETED'
).create();
