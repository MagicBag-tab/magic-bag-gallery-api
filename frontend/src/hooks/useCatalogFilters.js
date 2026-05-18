import { useReducer, useMemo } from 'react';

export const initialState = {
  pinturas: [],
  search: '',
  filterExclusiva: 'all',
  loading: true,
  error: null,
};

export function catalogReducer(state, action) {
  switch (action.type) {
    case 'SET_PINTURAS':
      return { ...state, pinturas: action.payload, loading: false, error: null };
    case 'SET_LOADING':
      return { ...state, loading: action.payload };
    case 'SET_ERROR':
      return { ...state, error: action.payload, loading: false };
    case 'SET_SEARCH':
      return { ...state, search: action.payload };
    case 'SET_FILTER_EXCLUSIVA':
      return { ...state, filterExclusiva: action.payload };
    case 'RESET_FILTERS':
      return { ...state, search: '', filterExclusiva: 'all' };
    default:
      return state;
  }
}

export function useCatalogFilters() {
  const [state, dispatch] = useReducer(catalogReducer, initialState);

  const filtered = useMemo(() => {
    let result = state.pinturas;

    if (state.search) {
      const q = state.search.toLowerCase();
      result = result.filter(
        (p) =>
          p.titulo.toLowerCase().includes(q) ||
          p.artista.toLowerCase().includes(q) ||
          p.coleccion?.toLowerCase().includes(q)
      );
    }

    if (state.filterExclusiva === 'si') {
      result = result.filter((p) => p.exclusiva);
    } else if (state.filterExclusiva === 'no') {
      result = result.filter((p) => !p.exclusiva);
    }

    return result;
  }, [state.pinturas, state.search, state.filterExclusiva]);

  const setPinturas = (pinturas) =>
    dispatch({ type: 'SET_PINTURAS', payload: pinturas });
  const setLoading = (val) =>
    dispatch({ type: 'SET_LOADING', payload: val });
  const setError = (msg) =>
    dispatch({ type: 'SET_ERROR', payload: msg });
  const setSearch = (val) =>
    dispatch({ type: 'SET_SEARCH', payload: val });
  const setFilterExclusiva = (val) =>
    dispatch({ type: 'SET_FILTER_EXCLUSIVA', payload: val });
  const resetFilters = () =>
    dispatch({ type: 'RESET_FILTERS' });

  return {
    ...state,
    filtered,
    setPinturas,
    setLoading,
    setError,
    setSearch,
    setFilterExclusiva,
    resetFilters,
  };
}