# 容器化部署Kafka

Strimzi项目实现了通过Kubernetes Operator方式统一部署，管理，维护容器化Kafka的功能。Strimzi主要功能包括了部署运行Kafka集群，部署Kafka组件，维护Kafka访问方式。安全加固，Broker管理和Kafka升级。除此之外，Strimzi还支持通过Kubernets CRD的方式管理Kafka的Topic与User。

### Strimzi准备

使用版本：0.20.0

Github：https://github.com/strimzi/strimzi-kafka-operator/releases

### 使用步骤

以下步骤需要kubectl使用的kubeconfig的账号有Kuberntes集群的admin权限。

#### 创建Kafka Operator的命名空间

```shell
kubectl create ns kafka
```

该命名空间用于运行Kafka Operator；Kafka Operator用于监听给定的Kafka集群运行使用的命名空间，并根据内部资源变化进行对应的管理性操作。

该namespace名称可修改，如下：

```shell
sed -i 's/namespace: .*/namespace: kafka/' install/cluster-operator/*RoleBinding*.yaml
```

修改以上指令的`kafka`即可，即：

```shell
sed -i 's/namespace: .*/namespace: strimzi/' install/cluster-operator/*RoleBinding*.yaml
```

#### 创建用于部署Kafka Cluster的命名空间

```shell
kubectl create ns my-kafka-project
```

同时修改` install/cluster-operator/060-Deployment-strimzi-cluster-operator.yaml `下的 `STRIMZI_NAMESPACE`  变量：

```yaml
# ...
env:
- name: STRIMZI_NAMESPACE
  value: my-kafka-project
# ...
```

#### 部署Strimzi相关的CRD和RBAC

```shell
kubectl apply -f install/cluster-operator/ -n kafka
```

#### 为Kafka Operator添加访问Kafka Cluster的RBAC

```shell
kubectl apply -f install/cluster-operator/020-RoleBinding-strimzi-cluster-operator.yaml -n my-kafka-project
kubectl apply -f install/cluster-operator/032-RoleBinding-strimzi-cluster-operator-topic-operator-delegation.yaml -n my-kafka-project
kubectl apply -f install/cluster-operator/031-RoleBinding-strimzi-cluster-operator-entity-operator-delegation.yaml -n my-kafka-project
```

#### 创建Kafka Cluster

```yaml
apiVersion: kafka.strimzi.io/v1beta1
kind: Kafka
metadata:
  name: my-cluster
  namespace: my-kafka-project
spec:
  kafka:
    replicas: 1	# kafka副本数
    resources:	# kafka占用资源
   	  requests:
        memory: 500Mi
    	cpu: 0.5
      limits:
        memory: 1Gi
    	cpu: 1
    listeners:
      - name: plain
        port: 9092
        type: internal
        tls: false
      - name: tls
        port: 9093
        type: internal
        tls: true
        authentication:
          type: tls
      - name: external	# 对外通过nodeport访问
        port: 9094
        type: nodeport
        tls: false
    template:
      pod:
        securityContext:
          runAsUser: 0
          fsGroup: 0
    storage:
      type: jbod
      volumes:
      - id: 0
        type: persistent-claim
        size: 100Gi
        deleteClaim: false
        class: cephfs
    config:	# kafka相关配置
      offsets.topic.replication.factor: 1
      transaction.state.log.replication.factor: 1
      transaction.state.log.min.isr: 1
  zookeeper:
    replicas: 1	# zookeeper副本数
   	resources:	# zookeeper占用资源
   	  requests:
        memory: 500Mi
        cpu: 0.5
     limits:
        memory: 1Gi
         cpu: 1
    storage:
      type: persistent-claim
      size: 100Gi
      deleteClaim: false
      class: cephfs
    template:
      pod:
        securityContext:
          runAsUser: 0
          fsGroup: 0
  entityOperator:
    topicOperator: {}
    userOperator: {}
```

注：对不同的StorageClass，需要指定适当的容器权限。

#### 创建Kafka Topic

```yaml
apiVersion: kafka.strimzi.io/v1beta1
kind: KafkaTopic
metadata:
  name: my-topic	# topic名称
  namespace: my-kafka-project
  labels:
    strimzi.io/cluster: "my-cluster"	# 对应Kafka集群
spec:
  partitions: 1	# 分区数量
  replicas: 1	# topic中各个分区的副本数量
```

部署完成后，可通过`IP:NodePort`的方式访问已经创建的Kafka Topic。

#### 添加Kafka User

如要添加Kafka用户，示例配置如下：

```yaml
apiVersion: kafka.strimzi.io/v1beta1
kind: KafkaUser
metadata:
  name: my-user
  labels:
    strimzi.io/cluster: my-cluster
spec:
  authentication:
    type: tls
  authorization:
    type: simple
    acls:
      # Example consumer Acls for topic my-topic using consumer group my-group
      - resource:
          type: topic
          name: my-topic
          patternType: literal
        operation: Read
        host: "*"
      - resource:
          type: topic
          name: my-topic
          patternType: literal
        operation: Describe
        host: "*"
      - resource:
          type: group
          name: my-group
          patternType: literal
        operation: Read
        host: "*"
      # Example Producer Acls for topic my-topic
      - resource:
          type: topic
          name: my-topic
          patternType: literal
        operation: Write
        host: "*"
      - resource:
          type: topic
          name: my-topic
          patternType: literal
        operation: Create
        host: "*"
      - resource:
          type: topic
          name: my-topic
          patternType: literal
        operation: Describe
        host: "*"
```

