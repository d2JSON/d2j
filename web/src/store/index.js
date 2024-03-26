const getDatabaseKeyToLocalStorage = (databaseKey) => {
  return localStorage.getItem('databaseKey', databaseKey)
}

const saveDatabaseKeyToLocalStorage = (databaseKey) => {
  localStorage.setItem('databaseKey', databaseKey)
}

export {
  getDatabaseKeyToLocalStorage,
  saveDatabaseKeyToLocalStorage
}
