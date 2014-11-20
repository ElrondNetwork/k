// Copyright (c) 2014 K Team. All Rights Reserved.

package org.kframework.kore

import collection.{ AbstractSeq, LinearSeq, LinearSeqOptimized, Seq, generic, mutable }
import collection.JavaConverters._
import java.util.stream.StreamSupport

abstract class KList extends KListLike[KList] with KCollection with KListMatcher {
  type ThisK = KList
  def copy(l: LinearSeq[K]) = KList(l: _*)
  def copy(l: Iterable[K]) = copy(l.toList)

  override def toString = this.mkString(", ")
}

final case object EmptyKList extends KList with Serializable {
  override def isEmpty = true
  override def head: Nothing = throw new NoSuchElementException("head of empty klist")
  override def tail: KList = throw new UnsupportedOperationException("tail of empty klist")
  override def equals(that: Any) = that match {
    case that: scala.collection.GenSeq[_] => that.isEmpty
    case _ => false
  }
}

final case class ConsKList(override val head: K, override val tail: KList) extends KList {
  override def isEmpty = false
}

object KList extends CanBuildKListLike[KList] {
  def apply(l: K*): KList =
    l.foldRight(EmptyKList: KList) {
      case (KApply(KLabel("KList"), h, _), l: KList) => KList((h ++ l).toSeq: _*)
      case (h: K, l: KList) => new ConsKList(h, l)
    }

  implicit def inject(k: K): KList = KList(k)
  implicit def seqOfKtoKList(s: Seq[K]) = KList(s: _*)
  def fromJava(l: Array[K]) = apply(l: _*)
  val fromList = seqOfKtoKList _

  def unapplySeq(l: KCollection): Option[Seq[K]] = Some(l.toSeq)
}
