import { reducerFactory } from '../../../../core/redux';
import {
  panelEditorCleanUp,
  panelEditorInitCompleted,
  panelEditorTabTexts,
  PanelEditorTab,
  PanelEditorTabIds,
} from './types';

export const getPanelEditorTab = (tabId: PanelEditorTabIds): PanelEditorTab => {
  return {
    id: tabId,
    text: panelEditorTabTexts[tabId],
  };
};

export interface PanelEditorState {
  activeTab: PanelEditorTabIds;
  tabs: PanelEditorTab[];
}

export const initialState: PanelEditorState = {
  activeTab: null,
  tabs: [],
};

export const panelEditorReducer = reducerFactory<PanelEditorState>(initialState)
  .addMapper({
    filter: panelEditorInitCompleted,
    mapper: (state, action): PanelEditorState => {
      const { activeTab, tabs } = action.payload;
      return {
        ...state,
        activeTab,
        tabs,
      };
    },
  })
  .addMapper({
    filter: panelEditorCleanUp,
    mapper: (): PanelEditorState => initialState,
  })
  .create();
