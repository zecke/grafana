import { actionCreatorFactory } from '../../../../core/redux';

export interface PanelEditorTab {
  id: string;
  text: string;
}

export enum PanelEditorTabIds {
  Queries = 'queries',
  Visualization = 'visualization',
  Advanced = 'advanced',
  Alert = 'alert',
}

export const panelEditorTabTexts = {
  [PanelEditorTabIds.Queries]: 'Queries',
  [PanelEditorTabIds.Visualization]: 'Visualization',
  [PanelEditorTabIds.Advanced]: 'General',
  [PanelEditorTabIds.Alert]: 'Alert',
};

export interface PanelEditorInitCompleted {
  activeTab: PanelEditorTabIds;
  tabs: PanelEditorTab[];
}

export const panelEditorCleanUp = actionCreatorFactory('PANEL_EDITOR_CLEAN_UP').create();

export const panelEditorInitCompleted = actionCreatorFactory<PanelEditorInitCompleted>(
  'PANEL_EDITOR_INIT_COMPLETED'
).create();
