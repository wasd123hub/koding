config =
  _selfConfig :
    defaultConfig : "dev"
  kite :
    applications :
      name              : 1
      pidPath           : 1
      logFile           : 1
      amqp            :
        host            : 1
        login           : 1
        password        : 1
        heartbeat       : 1
      apiUri            : 1
      usersPath         : 1
      vhostDir          : 1
      defaultDomain     : 1
      minAllowedUid     : 1
      debugApi          : 1
    databases    :
      name              : 1
      pidPath           : 1
      logFile           : 1
      port              : 1
      amqp                  :
        host                : 1
        login               : 1
        password            : 1
        heartbeat           : 1
      apiUri            : 1
      mysql             :
        usersPath       : 1
        backupDir       : 1
        databases       :
          mysql         : [{ host : 1, user : 1, password:1}]
      mongo             :
        databases       :
          mongodb       : [{ host : 1, user : 1, password:1}]
    sharedHosting :
      name                  : 1
      numberOfWorkers       : 1
      pidPath               : 1
      logFile               : 1
      amqp                  :
        host                : 1
        login               : 1
        password            : 1
        heartbeat           : 1
      apiUri                : 1
      usersPath             : 1
      vhostDir              : 1
      suspendDir            : 1
      defaultVhostFiles     : 1
      freeUsersGroup        : 1
      liteSpeedUser         : 1
      defaultDomain         : 1
      minAllowedUid         : 1
      debugApi              : 1
      processBaseDir        : 1
      cagefsctl             : 1
      baseMountDir          : 1
      maxAllowedRemotes     : 1
      usersMountsFile       : 1
      encryptKey            : 1
      ftpfs                 :
        curlftpfs           : 1
        opts                : 1
      sshfs                 :
        sshfscmd            : 1
        opts                : 1
        optsWithKey         : 1
      lsws                  :
        baseDir             : 1
        controllerPath      : 1
        lsMasterConfig      : 1
        configFilePath      : 1
        minRestartInterval  : 1
      ldap                  :
        ldapUrl             : 1
        rootUser            : 1
        rootPass            : 1
        groupDN             : 1
        userDN              : 1
        freeUID             : 1
        freeGroup           : 1
      FileSharing           :
        baseSharedDir       : 1
        baseDir             : 1
        setfacl             : 1
  main :
    haproxy       :
      webPort     : 1
    aws           :
      key         : 1
      secret      : 1
    uri           :
      address     : 1
    projectRoot   : 1
    version       : 1
    webserver     :
      login       : 1
      port        : []
      clusterSize : 1
      queueName   : 1
      watch       : 1
    sourceServer  :
      enabled     : 1
      port        : 1
    mongo         : 1
    neo4j         :
      enabled     : 1
      host        : 1
      port        : 1
    runGoBroker   : 1
    compileGo     : 1
    buildClient   : 1
    misc          :
      claimGlobalNamesForUsers: 1
      updateAllSlugs : 1
      debugConnectionErrors: 1
    uploads       :
      enableStreamingUploads: 1
      distribution: 1
      s3          :
        awsAccountId        : 1
        awsAccessKeyId      : 1
        awsSecretAccessKey  : 1
        bucket              : 1
    loggr:
      push: 1
      url: 1
      apiKey: 1
    librato :
      push      : 1
      email     : 1
      token     : 1
      interval  : 1
    goConfig:
      HomePrefix   : 1
      UseLVE       : 1
    bitly :
      username  : 1
      apiKey    : 1
    authWorker    :
      login           : 1
      queueName       : 1
      authResourceName: 1
      numberOfWorkers : 1
      watch           : 1
    social        :
      login       : 1
      numberOfWorkers: 1
      watch       : 1
      queueName   : 1
    cacheWorker       :
      login           : 1
      watch           : 1
      queueName       : 1
      run             : 1
    feeder        :
      queueName   : 1
      exchangePrefix: 1
      numberOfWorkers: 1
    presence      :
      exchange    : 1
    client        :
      version     : 1
      watch       : 1
      includesPath: 1
      websitePath : 1
      js          : 1
      css         : 1
      indexMaster : 1
      index       : 1
      useStaticFileServer: 1
      staticFilesBaseUrl: 1
      runtimeOptions:
        logToExternal : 1
        resourceName  : 1
        suppressLogs  : 1
        version       : 1
        mainUri       : 1
        broker        :
          sockJS      : 1
        apiUri        : 1
        appsUri       : 1
        sourceUri     : 1
    mq            :
      host        : 1
      port        : 1
      apiPort     : 1
      login       : 1
      componentUser: 1
      password    : 1
      heartbeat   : 1
      vhost       : 1
    broker        :
      ip          : 1
      port        : 1
      certFile    : 1
      keyFile     : 1
    kites:
      disconnectTimeout: 1
      vhost       : 1
    email         :
      host        : 1
      protocol    : 1
      defaultFromAddress: 1
    emailWorker   :
      cronInstant : 1
      cronDaily   : 1
      run         : 1
      defaultRecepient : 1
    guests        :
      poolSize        : 1
      batchSize       : 1
      cleanupCron     : 1
    pidFile       : 1


module.exports = config
