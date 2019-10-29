import { reducerTester } from '../../../../../test/core/redux/reducerTester';
import { initialState, panelEditorReducer, PanelEditorTabIds, PanelEditorTab, getPanelEditorTab } from './reducers';
import { panelEditorCleanUp } from './actions';

describe('panelEditorReducer', () => {
  describe('when panelEditorCleanUp is dispatched', () => {
    it('then state should be intialState', () => {
      const tabs: PanelEditorTab[] = [
        getPanelEditorTab(PanelEditorTabIds.Queries),
        getPanelEditorTab(PanelEditorTabIds.Visualization),
        getPanelEditorTab(PanelEditorTabIds.Advanced),
      ];
      reducerTester()
        .givenReducer(panelEditorReducer, { tabs })
        .whenActionIsDispatched(panelEditorCleanUp())
        .thenStateShouldEqual(initialState);
    });
  });
});
